package go_watchman

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sotera/go_watchman/loogo"
	"log"
	"os"
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

type annotationOptions struct {
	startTime        string
	endTime          string
	apiRoot          string
	annotationType   string
	annotation_types []string
	fetcher          fetcher
	pagerFactory     LoogoPagerFactory
}

type fetcher interface {
	fetch(options annotationOptions) ([]Annotation, error)
}

type LoogoPagerFactory interface {
	generate(params loogo.NewPagerParams) (loogo.PagerInterface, error)
}

func main() {
	options := annotationOptions{}
	options.startTime = *flag.String("start-time-ms", "", "start time in millis")
	options.endTime = *flag.String("end-time-ms", "", "end time in millis")
	options.annotationType = *flag.String("annotation-type", "", "the type of annotation to process")
	flag.Parse()

	options.apiRoot = os.Getenv("ANNOTATION_API_ROOT")
	if options.apiRoot == "" {
		options.apiRoot = "http://dev-qcr-io-services-qntfy-annotation-api.traefik.dsra.local:31888/v1/annotations"
	}

	options.annotation_types = []string{"name", "relevant"}
	options.fetcher = annotation_fetcher{}
	options.pagerFactory = pager_factory{}
	err := process_annotation_types(options)
	if err != nil {
		log.Fatal(fmt.Println(err))
	}
}

func process_annotation_types(options annotationOptions) error {
	if options.apiRoot == "" || options.endTime == "" || options.startTime == "" {
		return errors.New("invalid options")
	}
	if options.fetcher == nil {
		return errors.New("fetcher instance was nil, please provide a fetcher")
	}
	if options.annotation_types == nil || len(options.annotation_types) == 0 {
		return errors.New("no annotation types to process")
	}

	for i := 0; i < len(options.annotation_types); i++ {
		options.annotationType = options.annotation_types[i]
		annotations, err := fetch_annotations(options)
		if err != nil {
			return err
		}
		err = process_annotations(annotations, options.pagerFactory)
		if err != nil {
			return err
		}
	}

	return nil
}

func fetch_annotations(options annotationOptions) ([]Annotation, error) {
	annotations, err := options.fetcher.fetch(options)
	if err != nil {
		return nil, err
	}

	fmt.Println("annotations:", len(annotations))
	return annotations, nil
}

func parse_annotation_id(annotation_id string) (campaign string, event_id string) {
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

func process_annotations(annotations []Annotation, pagerFactory LoogoPagerFactory) error {
	fmt.Println("annotations:", len(annotations))
	var wg sync.WaitGroup
	for i := 0; i < len(annotations); i++ {
		annotation := annotations[i]
		annotation.CampaignId, annotation.EventId = parse_annotation_id(annotation.Object_id)

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

		pager, err := pagerFactory.generate(loogo.NewPagerParams{
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
			go create_annotation(&wg, annotation, page[0])
		} else {
			go update_annotation(&wg, annotation, page[0])
		}
	}

	wg.Wait()
	return nil
}

func create_annotation(wg *sync.WaitGroup, annotation Annotation, doc loogo.Doc) {
	defer wg.Done()

	//create new annotation
}

func update_annotation(wg *sync.WaitGroup, annotation Annotation, doc loogo.Doc) {
	defer wg.Done()

	//update annotation
}
