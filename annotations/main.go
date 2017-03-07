package go_watchman

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"github.com/sotera/go_watchman/loogo"
)

type Annotation struct {
	Object_id       string
	Reference_id    string
	Annotation_type string `json:"type"`
	Value           string
	Annotator       string
}

type annotationOptions struct {
	startTime      string
	endTime        string
	apiRoot        string
	annotationType string
}

type fetcher interface {
	fetch(options annotationOptions) ([]Annotation, error)
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

	annotation_types := []string{"name", "relevant"}
	fetcher := annotation_fetcher{}

	err := process_annotation_types(options, annotation_types, fetcher)
	if err != nil {
		log.Fatal(fmt.Println(err))
	}
}

func process_annotation_types(options annotationOptions, annotation_types []string, fetcher fetcher) error {

	if fetcher == nil {
		return errors.New("fetcher instance was nil, please provide a fetcher")
	}
	if options.apiRoot == "" || options.endTime == "" || options.startTime == "" {
		return errors.New("invalid options")
	}
	if annotation_types == nil || len(annotation_types) == 0 {
		return errors.New("no annotation types to process")
	}

	for i := 0; i < len(annotation_types); i++ {
		options.annotationType = annotation_types[i]
		annotations, err := fetch_annotations(options, fetcher)
		if err != nil {
			return err
		}
		err = process_annotations(annotations)
		if err != nil {
			return err
		}
	}

	return nil
}

func fetch_annotations(options annotationOptions, fetcher fetcher) ([]Annotation, error) {
	annotations, err := fetcher.fetch(options)
	if err != nil {
		return nil, err
	}

	fmt.Println("annotations:", len(annotations))
	return annotations, nil
}

func parse_annotation_id(annotation_id string) (campaign string, event_id string) {
	tokens := strings.Split(annotation_id, ":")
	campaign = tokens[1]
	event_id = tokens[2]
	return
}

func process_annotations(annotations []Annotation) error {
	fmt.Println("annotations:", len(annotations))
	var wg sync.WaitGroup
	for i := 0; i < len(annotations); i++ {
		wg.Add(1)
		annotation := annotations[i]
		go update_event(&wg, annotation)
	}

	wg.Wait()
	return nil
}

func update_event(wg *sync.WaitGroup, annotation Annotation) {
	defer wg.Done()
	campaign, event_id := parse_annotation_id(annotation.Object_id)
	fmt.Printf("campaign: %v event_id: %v", campaign, event_id)

	p1 := loogo.QueryParam{
		QueryType: "Eq",
		Field:     "id",
		Values:    []string{"new"},
	}
	loogo.Eq(p1, false)
}
