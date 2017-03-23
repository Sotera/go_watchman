package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Sotera/go_watchman/annotations"
)

func main() {

	startTime := flag.String("start-time-ms", "", "start time in millis")
	endTime := flag.String("end-time-ms", "", "end time in millis")
	flag.Parse()

	options := annotations.AnnotationOptions{}
	options.StartTime = *startTime
	options.EndTime = *endTime

	options.AnnotationAPIRoot = os.Getenv("ANNOTATION_API_ROOT")
	if options.AnnotationAPIRoot == "" {
		options.AnnotationAPIRoot = "http://dev-qcr-io-services-qntfy-annotation-api.traefik.dsra.local:31888/v1/annotation"
	} else {
		options.AnnotationAPIRoot = strings.TrimRight(options.AnnotationAPIRoot, "/")
	}

	options.AnnotationRefID = os.Getenv("ANNOTATION_REF_ID")
	if options.AnnotationRefID == "" {
		options.AnnotationRefID = "qcr.app.dev"
	}

	options.APIRoot = os.Getenv("API_ROOT")
	if options.APIRoot == "" {
		options.APIRoot = "http://localhost:3003/api"
	} else {
		options.APIRoot = strings.TrimRight(options.APIRoot, "/")
	}

	options.AnnotationTypes = []string{"label", "relevant"}
	options.Fetcher = annotations.AnnotationFetcher{}

	err := annotations.ProcessAnnotationTypes(options)
	if err != nil {
		log.Println(fmt.Println(err))
	}
}
