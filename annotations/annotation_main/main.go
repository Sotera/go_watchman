package main

import (
	"flag"
	"fmt"
	"github.com/Sotera/go_watchman/annotations"
	"log"
	"os"
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

	options.Annotation_types = []string{"name", "relevant"}
	options.Fetcher = annotations.AnnotationFetcher{}
	options.PagerFactory = annotations.PagerFactory{}
	err := annotations.ProcessAnnotationTypes(options)
	if err != nil {
		log.Fatal(fmt.Println(err))
	}
}
