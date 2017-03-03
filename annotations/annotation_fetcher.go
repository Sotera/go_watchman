package go_watchman

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)


type annotation_fetcher struct {
}

func (af annotation_fetcher) fetch(options annotationOptions) ([]Annotation, error) {

	flag.Parse()

	res, err := http.Get(options.apiRoot + "/type/" + options.annotationType)

	if err != nil {
		log.Fatal(fmt.Println(err))
		return nil, err
	}

	annotations := make([]Annotation, 0)

	decodeErr := json.NewDecoder(res.Body).Decode(&annotations)
	if decodeErr != nil {
		log.Fatal(fmt.Println(decodeErr))
		return nil, decodeErr
	}

	fmt.Println("annotations:", len(annotations))
	return annotations, nil
}
