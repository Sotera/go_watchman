package annotations

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Fetcher interface {
	Fetch(options AnnotationOptions) ([]Annotation, error)
}

type AnnotationFetcher struct {
}

func getDateRange(options AnnotationOptions) (string, string) {
	defaultTimeStr := fmt.Sprintf("%d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day())

	start, err1 := strconv.ParseInt(options.StartTime, 10, 64)
	if err1 != nil {
		log.Println("bad or no start time provided...returning default time")
		return defaultTimeStr, defaultTimeStr
	}
	startTime := time.Unix(0, start*int64(time.Millisecond))
	startStr := fmt.Sprintf("%d-%02d-%02d", startTime.Year(), startTime.Month(), startTime.Day())

	end, err2 := strconv.ParseInt(options.EndTime, 10, 64)
	if err2 != nil {
		log.Println("bad or no end time provided...returning default time")
		return defaultTimeStr, defaultTimeStr
	}
	endTime := time.Unix(0, end*int64(time.Millisecond))
	endStr := fmt.Sprintf("%d-%02d-%02d", endTime.Year(), endTime.Month(), endTime.Day())

	return startStr, endStr
}

func (af AnnotationFetcher) Fetch(options AnnotationOptions) ([]Annotation, error) {

	startStr, endStr := getDateRange(options)
	url := fmt.Sprintf("%s/type/%s/?from_date=%s&to_date=%s", options.AnnotationAPIRoot, options.AnnotationType, startStr, endStr)
	println(url)
	res, err := http.Get(url)

	if err != nil {
		log.Println(fmt.Println(err))
		return nil, err
	}

	annotations := make([]Annotation, 0)

	decodeErr := json.NewDecoder(res.Body).Decode(&annotations)
	if decodeErr != nil {
		log.Println(fmt.Println(decodeErr))
		return nil, decodeErr
	}
	//fake for now
	/*annotations := []Annotation{{
	ObjectID:       "smevent:MyTestCampaign:06c1909d-0ce3-4df3-86fd-2104c10a8581",
	ReferenceID:    "qcr.app.dev",
	AnnotationType: "name",
	Value:          "test3",
	Annotator:      "alex"}}*/
	fmt.Println("annotations:", len(annotations))
	return annotations, nil
}
