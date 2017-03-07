package annotations

import (
	"errors"
	"fmt"
	"github.com/sotera/go_watchman/loogo"
	"strings"
	"sync"
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

type AnnotationOptions struct {
	StartTime         string
	EndTime           string
	ApiRoot           string
	AnnotationApiRoot string
	AnnotationType    string
	Annotation_types  []string
	Fetcher           Fetcher
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
		err = ProcessAnnotations(annotations, options.PagerFactory)
		if err != nil {
			return err
		}
	}

	return nil
}

func FetchAnnotations(options AnnotationOptions) ([]Annotation, error) {
	annotations, err := options.Fetcher.Fetch(options)
	if err != nil {
		return nil, err
	}

	fmt.Println("annotations:", len(annotations))
	return annotations, nil
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

func ProcessAnnotations(annotations []Annotation, pagerFactory LoogoPagerFactory) error {
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

		pager, err := pagerFactory.Generate(loogo.NewPagerParams{
			URL:      "http://localhost/api/events",
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
			go CreateAnnotation(&wg, annotation, page[0])
		} else {
			go UpdateAnnotation(&wg, annotation, page[0])
		}
	}

	wg.Wait()
	return nil
}

func CreateAnnotation(wg *sync.WaitGroup, annotation Annotation, doc loogo.Doc) {
	defer wg.Done()

	//create new annotation
}

func UpdateAnnotation(wg *sync.WaitGroup, annotation Annotation, doc loogo.Doc) {
	defer wg.Done()

	//update annotation
}
