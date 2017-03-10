package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Sotera/go_watchman/annotations"
	"github.com/Sotera/go_watchman/loogo"
)

func main() {

	startTime := flag.String("start-time-ms", "", "start time in millis")
	endTime := flag.String("end-time-ms", "", "end time in millis")
	flag.Parse()

	options := annotations.AnnotationOptions{}
	options.StartTime = *startTime
	options.EndTime = *endTime

	options.AnnotationApiRoot = os.Getenv("ANNOTATION_API_ROOT")
	if options.AnnotationApiRoot == "" {
		options.AnnotationApiRoot = "http://dev-qcr-io-services-qntfy-annotation-api.traefik.dsra.local:31888/v1/annotations"
	}

	options.ApiRoot = os.Getenv("API_ROOT")
	if options.ApiRoot == "" {
		options.ApiRoot = "http://localhost:3003/api"
	}

	parser := loogo.HTTPRequestParser{
		Client: &loogo.HTTPClient{},
	}

	options.Annotation_types = []string{"name", "relevant"}
	options.Fetcher = annotations.AnnotationFetcher{}
	options.PagerFactory = annotations.PagerFactory{}
	options.Parser = &parser

	err := annotations.ProcessAnnotationTypes(options)
	if err != nil {
		log.Fatal(fmt.Println(err))
	}
}
