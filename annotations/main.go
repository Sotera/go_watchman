package go_watchman

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

type annotationOptions struct {
	startTime      string
	endTime        string
	apiRoot        string
	annotationType string
}

type fetcher interface {
	fetch(options annotationOptions) (*[]Annotation, error)
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

func fetch_annotations(options annotationOptions, fetcher fetcher) (*[]Annotation, error) {
	annotations, err := fetcher.fetch(options)
	if err != nil {
		return nil, err
	}

	fmt.Println("annotations:", len(*annotations))
	return annotations, nil
}

func process_annotations(annotations *[]Annotation) error {
	fmt.Println("annotations:", len(*annotations))

	return nil
}
