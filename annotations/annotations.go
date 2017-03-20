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

type Annotation struct {
	ObjectID       string `json:"Object_id"`
	ReferenceID    string `json:"Reference_id"`
	AnnotationType string `json:"type"`
	Value          string `json:"Value"`
	Annotator      string `json:"Annotator"`
	EventID        string `json:"EventId"`
	CampaignID     string `json:"CampaignId"`
}

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
	ParserFactory     LoogoParserFactory
	PagerFactory      LoogoPagerFactory
}

func ProcessAnnotationTypes(options AnnotationOptions) error {
	if options.Fetcher == nil {
		return errors.New("fetcher instance was nil, please provide a fetcher")
	}
	if options.AnnotationTypes == nil || len(options.AnnotationTypes) == 0 {
		return errors.New("no annotation types to process")
	}

	for i := 0; i < len(options.AnnotationTypes); i++ {
		options.AnnotationType = options.AnnotationTypes[i]
		annotations, err := FetchAnnotations(options)
		if err != nil {
			log.Println(err)
			continue
		}
		err = ProcessAnnotations(annotations, options)
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

	fmt.Println(fmt.Sprintf("fetcher provided %d annotations", len(results)))
	return results, nil
}

func ParseAnnotationID(annotation_id string) (campaign string, event_id string) {
	tokens := strings.Split(annotation_id, ":")
	campaign, event_id = "", ""
	if tokens == nil || len(tokens) <= 1 {
		return
	}
	if len(tokens) > 1 {
		campaign = tokens[1]
	}
	if len(tokens) > 2 {
		event_id = tokens[2]
	}
	return
}

func ProcessAnnotations(annotations []Annotation, options AnnotationOptions) error {
	var wg sync.WaitGroup
	if options.PagerFactory == nil {
		return errors.New("options.PagerFactory was nil")
	}

	if options.ParserFactory == nil {
		return errors.New("options.ParserFactory was nil")
	}

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

		pager, err := options.PagerFactory.Generate(loogo.NewPagerParams{
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
			go CreateAnnotation(&wg, options, annotation)
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

			go UpdateAnnotation(&wg, options, annotation, model)
		}
	}

	wg.Wait()
	return nil
}

func GetEvent(eventChannel chan loogo.Doc, options AnnotationOptions, annotation Annotation) {
	url := fmt.Sprintf("%s/events/%s", options.APIRoot, annotation.EventID)
	println(url)

	event := loogo.Doc{}

	parser := options.ParserFactory.Generate()
	err := parser.NewRequest(loogo.NewRequestParams{URL: url}, &event)
	if err != nil {
		log.Println(err)
		eventChannel <- nil
		return
	}
	eventChannel <- event
}

func CreateAnnotation(wg *sync.WaitGroup, options AnnotationOptions, annotation Annotation) {
	defer wg.Done()
	model := AnnotationModel{}
	model.Campaign = annotation.CampaignID
	model.Event = annotation.EventID
	println(annotation.EventID)

	eventChannel := make(chan loogo.Doc)
	go GetEvent(eventChannel, options, annotation)
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
	options.ParserFactory.Generate().NewRequest(params, doc)
}

func UpdateAnnotation(wg *sync.WaitGroup, options AnnotationOptions, annotation Annotation, model AnnotationModel) {
	defer wg.Done()
	var m map[string]int
	m = make(map[string]int)
	m["a"] = 5
	m["b"] = 6
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
	options.ParserFactory.Generate().NewRequest(params, doc)
}
