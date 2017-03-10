package annotations

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/Sotera/go_watchman/loogo"
)

type Annotation struct {
	Object_id       string
	Reference_id    string
	Annotation_type string `json:"type"`
	Value           string
	Annotator       string
	EventId         string
	CampaignId      string
}

type AnnotationModel struct {
	Campaign string   `json:"campaign,omitempty"`
	Event    string   `json:"event,omitempty"`
	Features []string `json:"features,omitempty"`
	Name     string   `json:"name,omitempty"`
	Relevant bool     `json:"relevant,omitempty"`
	ID       string   `json:"id,omitempty"`
}

type AnnotationOptions struct {
	StartTime         string
	EndTime           string
	ApiRoot           string
	AnnotationApiRoot string
	AnnotationType    string
	Annotation_types  []string
	Fetcher           Fetcher
	Parser            loogo.RequestParser
	PagerFactory      LoogoPagerFactory
}

func ProcessAnnotationTypes(options AnnotationOptions) error {
	if options.AnnotationApiRoot == "" || options.EndTime == "" || options.StartTime == "" {
		return errors.New("invalid options")
	}
	if options.Fetcher == nil {
		return errors.New("fetcher instance was nil, please provide a fetcher")
	}
	if options.Annotation_types == nil || len(options.Annotation_types) == 0 {
		return errors.New("no annotation types to process")
	}

	for i := 0; i < len(options.Annotation_types); i++ {
		options.AnnotationType = options.Annotation_types[i]
		annotations, err := FetchAnnotations(options)
		if err != nil {
			return err
		}
		err = ProcessAnnotations(annotations, options)
		if err != nil {
			return err
		}
	}

	return nil
}

func FetchAnnotations(options AnnotationOptions) ([]Annotation, error) {
	results, err := options.Fetcher.Fetch(options)
	if err != nil {
		return nil, err
	}

	fmt.Println("annotations:", len(results))
	return results, nil
}

func ParseAnnotationId(annotation_id string) (campaign string, event_id string) {
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
	fmt.Println("annotations:", len(annotations))
	var wg sync.WaitGroup
	for i := 0; i < len(annotations); i++ {
		annotation := annotations[i]
		annotation.CampaignId, annotation.EventId = ParseAnnotationId(annotation.Object_id)

		params := loogo.QueryParams{
			loogo.QueryParam{
				QueryType: "Eq",
				Field:     "event",
				Values:    []string{annotation.EventId},
			},
			loogo.QueryParam{
				QueryType: "Eq",
				Field:     "campaign",
				Values:    []string{annotation.CampaignId},
			},
		}

		pager, err := options.PagerFactory.Generate(loogo.NewPagerParams{
			URL:      "http://localhost:3003/api/annotations",
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
			err = CreateAnnotation(&wg, options, annotation)
		} else {
			model := AnnotationModel{}
			model.ID = page[0]["id"].(string)
			model.Campaign = page[0]["campaign"].(string)
			model.Event = page[0]["event"].(string)
			model.Name = page[0]["name"].(string)
			model.Relevant = page[0]["relevant"].(bool)
			err = UpdateAnnotation(&wg, options, annotation, model)
		}
	}

	wg.Wait()
	return nil
}

func CreateAnnotation(wg *sync.WaitGroup, options AnnotationOptions, annotation Annotation) error {
	defer wg.Done()
	model := AnnotationModel{}
	model.Campaign = annotation.CampaignId
	model.Event = annotation.EventId
	model.Features = []string{}
	if annotation.Annotation_type == "name" {
		model.Name = annotation.Value
	} else {
		model.Relevant, _ = strconv.ParseBool(annotation.Value)
	}
	doc := loogo.Doc{}
	bytes, err := json.Marshal(model)
	if err != nil {
		return err
	}
	params := loogo.NewRequestParams{
		URL:        options.ApiRoot + "/annotations",
		Body:       bytes,
		HTTPMethod: "POST",
	}
	options.Parser.NewRequest(params, doc)
	return nil
}

func UpdateAnnotation(wg *sync.WaitGroup, options AnnotationOptions, annotation Annotation, model AnnotationModel) error {
	defer wg.Done()
	var m map[string]int
	m = make(map[string]int)
	m["a"] = 5
	m["b"] = 6
	if annotation.Annotation_type == "name" {
		model.Name = annotation.Value
	} else {
		model.Relevant, _ = strconv.ParseBool(annotation.Value)
	}
	doc := loogo.Doc{}
	bytes, err := json.Marshal(model)
	if err != nil {
		return err
	}

	params := loogo.NewRequestParams{
		URL:        options.ApiRoot + "/annotations",
		Body:       bytes,
		HTTPMethod: "PUT",
	}
	options.Parser.NewRequest(params, doc)
	return nil
}
