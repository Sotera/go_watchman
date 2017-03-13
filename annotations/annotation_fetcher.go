package annotations

import (
	"flag"
	"fmt"
)

type Fetcher interface {
	Fetch(options AnnotationOptions) ([]Annotation, error)
}

type AnnotationFetcher struct {
}

func (af AnnotationFetcher) Fetch(options AnnotationOptions) ([]Annotation, error) {

	flag.Parse()

	/*res, err := http.Get(options.annotationApiRoot + "/type/" + options.annotationType)

	if err != nil {
		log.Fatal(fmt.Println(err))
		return nil, err
	}

	annotations := make([]Annotation, 0)


	decodeErr := json.NewDecoder(res.Body).Decode(&annotations)
	if decodeErr != nil {
		log.Fatal(fmt.Println(decodeErr))
		return nil, decodeErr
	}*/
	//fake for now
	annotations := []Annotation{{
		ObjectID:       "smevent:MyTestCampaign:06c1909d-0ce3-4df3-86fd-2104c10a8581",
		ReferenceID:    "qcr.app.dev",
		AnnotationType: "name",
		Value:          "test3",
		Annotator:      "alex"}}
	fmt.Println("annotations:", len(annotations))
	return annotations, nil
}
