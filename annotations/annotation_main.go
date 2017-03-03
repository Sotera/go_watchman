package go_watchman

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type AnnotationOptions struct {
	startTime      string
	endTime        string
	apiRoot        string
	annotationType string
}

func main() {

	options := AnnotationOptions{}
	options.startTime = *flag.String("start-time-ms", "", "start time in millis")
	options.endTime = *flag.String("end-time-ms", "", "end time in millis")
	options.annotationType = *flag.String("annotation-type", "", "the type of annotation to process")
	flag.Parse()

	options.apiRoot = os.Getenv("ANNOTATION_API_ROOT")
	if options.apiRoot == "" {
		options.apiRoot = "http://dev-qcr-io-services-qntfy-annotation-api.traefik.dsra.local:31888/v1/annotations"
	}

	annotations, err := Fetch_annotations(options)
	if err != nil {
		log.Fatal(fmt.Println(err))
	}
}
