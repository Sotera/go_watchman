package annotations

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/Sotera/go_watchman/loogo"
)

// Annotation models QCR's expected input object.
type Annotation struct {
	ObjectID       string `json:"Object_id"`
	ReferenceID    string `json:"Reference_id"`
	AnnotationType string `json:"type"`
	Value          string `json:"Value"`
	Annotator      string `json:"Annotator"`
	EventID        string `json:"EventId"`
	CampaignID     string `json:"CampaignId"`
}

// AnnotationModel models watchman event.
type AnnotationModel struct {
	Campaign string      `json:"campaign,omitempty"`
	Event    string      `json:"event,omitempty"`
	Features interface{} `json:"features,omitempty"`
	Name     string      `json:"name,omitempty"`
	Relevant bool        `json:"relevant,omitempty"`
	ID       string      `json:"id,omitempty"`
}

type AnnotationOptions struct {
	StartTime         string
	EndTime           string
	APIRoot           string
	AnnotationRefID   string
	AnnotationAPIRoot string
	AnnotationType    string
	AnnotationTypes   []string
	Fetcher           Fetcher
}

// AnnotationMaker takes care of creating annotations.
type AnnotationMaker struct {
	// we need a mockable instance generator, to create instances in a nested loop.
	// a 'normal' type won't work: we need to create many instances.
	// instead, use a func type to create new instances of a Pager interface.
	createPager   func(loogo.NewPagerParams) (loogo.PagerInterface, error)
	requestParser loogo.RequestParser
}

func ProcessAnnotationTypes(options AnnotationOptions) error {
	if options.Fetcher == nil {
		return errors.New("fetcher instance was nil, please provide a fetcher")
	}
	if options.AnnotationTypes == nil || len(options.AnnotationTypes) == 0 {
		return errors.New("no annotation types to process")
	}

	// tell annotationmaker how to create pagers.
	createPager := func(params loogo.NewPagerParams) (loogo.PagerInterface, error) {
		return loogo.NewPager(params)
	}
	//TODO: needs to be moved to calling func if we want to test this func.
	parser := &loogo.HTTPRequestParser{
		Client: &loogo.HTTPClient{},
	}

	am := &AnnotationMaker{createPager, parser}

	for i := 0; i < len(options.AnnotationTypes); i++ {
		options.AnnotationType = options.AnnotationTypes[i]
		annotations, err := FetchAnnotations(options)
		if err != nil {
			log.Println(err)
			continue
		}
		err = am.ProcessAnnotations(annotations, options)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	return nil
}

func FetchAnnotations(options AnnotationOptions) ([]Annotation, error) {
	results, err := options.Fetcher.Fetch(options)
	if err != nil {
		return nil, err
	}

	fmt.Printf("fetcher provided %d annotations\n", len(results))
	return results, nil
}

func ParseAnnotationID(annotationID string) (campaign string, eventID string) {
	tokens := strings.Split(annotationID, ":")
	campaign, eventID = "", ""
	if tokens == nil || len(tokens) <= 1 {
		return
	}
	if len(tokens) > 1 {
		campaign = tokens[1]
	}
	if len(tokens) > 2 {
		eventID = tokens[2]
	}
	return
}

func (am *AnnotationMaker) ProcessAnnotations(annotations []Annotation, options AnnotationOptions) error {
	var wg sync.WaitGroup

	for i := 0; i < len(annotations); i++ {
		annotation := annotations[i]
		annotation.CampaignID, annotation.EventID = ParseAnnotationID(annotation.ObjectID)

		params := loogo.QueryParams{
			loogo.QueryParam{
				QueryType: "Eq",
				Field:     "event",
				Values:    []string{annotation.EventID},
			},
			loogo.QueryParam{
				QueryType: "Eq",
				Field:     "campaign",
				Values:    []string{annotation.CampaignID},
			},
		}

		pager, err := am.createPager(loogo.NewPagerParams{
			URL:      options.APIRoot + "/annotations",
			Params:   params,
			PageSize: 1,
		})

		if err != nil {
			return err
		}

		page, err := pager.GetNext()
		if err != nil {
			return err
		}

		wg.Add(1)
		if len(page) < 1 {
			go am.CreateAnnotation(&wg, options, annotation)
		} else {
			model := AnnotationModel{}
			model.ID = page[0]["id"].(string)
			model.Campaign = page[0]["campaign"].(string)
			model.Event = page[0]["event"].(string)

			if name, ok := page[0]["name"]; ok {
				model.Name = name.(string)
			}
			if relevant, ok := page[0]["relevant"]; ok {
				model.Relevant = relevant.(bool)
			}

			go am.UpdateAnnotation(&wg, options, annotation, model)
		}
	}

	wg.Wait()
	return nil
}

func (am *AnnotationMaker) GetEvent(eventChannel chan loogo.Doc, options AnnotationOptions, annotation Annotation) {
	url := fmt.Sprintf("%s/events/%s", options.APIRoot, annotation.EventID)

	doc := loogo.Doc{}

	err := am.requestParser.NewRequest(loogo.NewRequestParams{URL: url}, &doc)
	if err != nil {
		log.Println(err)
		eventChannel <- nil
		return
	}
	eventChannel <- doc
}

func (am *AnnotationMaker) CreateAnnotation(wg *sync.WaitGroup, options AnnotationOptions, annotation Annotation) {
	defer wg.Done()
	model := AnnotationModel{}
	model.Campaign = annotation.CampaignID
	model.Event = annotation.EventID
	// println(annotation.EventID)

	eventChannel := make(chan loogo.Doc)
	go am.GetEvent(eventChannel, options, annotation)
	event := <-eventChannel
	if event == nil {
		log.Println("get event returned no event, event not found")
		return
	}

	model.Features = event["hashtags"]

	if annotation.AnnotationType == "label" {
		model.Name = annotation.Value
	} else {
		model.Relevant, _ = strconv.ParseBool(annotation.Value)
	}
	doc := loogo.Doc{}
	bytes, err := json.Marshal(model)
	if err != nil {
		log.Println(err)
		return
	}

	params := loogo.NewRequestParams{
		URL:        options.APIRoot + "/annotations",
		Body:       bytes,
		HTTPMethod: "POST",
	}
	err = am.requestParser.NewRequest(params, &doc)
	if err != nil {
		log.Println(err)
		return
	}
}

func (am *AnnotationMaker) UpdateAnnotation(wg *sync.WaitGroup, options AnnotationOptions, annotation Annotation, model AnnotationModel) {
	defer wg.Done()
	if annotation.AnnotationType == "label" {
		model.Name = annotation.Value
	} else {
		model.Relevant, _ = strconv.ParseBool(annotation.Value)
	}
	doc := loogo.Doc{}
	bytes, err := json.Marshal(model)
	if err != nil {
		log.Println(err)
		return
	}

	params := loogo.NewRequestParams{
		URL:        options.APIRoot + "/annotations",
		Body:       bytes,
		HTTPMethod: "PUT",
	}
	err = am.requestParser.NewRequest(params, &doc)
	if err != nil {
		log.Println(err)
		return
	}
}
