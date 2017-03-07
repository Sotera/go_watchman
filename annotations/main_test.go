package go_watchman

import (
	"github.com/sotera/go_watchman/loogo"
	"testing"
)

func TestFetchAnnotations(t *testing.T) {
	fetcher := mockFetcher{}

	pager_factory := mock_pager_factory{}

	annotation := Annotation{
		Object_id:       "smevent:campaignID:eventID",
		Reference_id:    "qcr.app.dev",
		Annotation_type: "name",
		Value:           "an event name",
		Annotator:       "alex"}

	fetcher.annotations = []Annotation{annotation}

	options := annotationOptions{
		startTime:        "",
		endTime:          "",
		apiRoot:          "",
		annotationType:   "",
		annotation_types: []string{"test"},
		fetcher:          fetcher,
		pagerFactory:     pager_factory,
	}

	val, err := fetch_annotations(options)

	if err != nil {
		t.Errorf("fetch failed with error: %v", err)
	}

	if len(val) != 1 {
		t.Errorf("annotation count expectd %v got %v", 1, len(val))
	}
}

func TestProcessAnnotationTypes(t *testing.T) {

	fetcher := mockFetcher{}

	annotation_types := []string{"name", "relevance"}

	options := annotationOptions{
		startTime:        "test",
		endTime:          "test",
		apiRoot:          "test",
		annotationType:   "test",
		fetcher:          nil,
		annotation_types: nil}

	err := process_annotation_types(options)
	if err == nil || err.Error() != "fetcher instance was nil, please provide a fetcher" {
		t.Errorf("fetcher instance was nil not caught: %v", err)
	}

	err = process_annotation_types(annotationOptions{})
	if err == nil || err.Error() != "invalid options" {
		t.Errorf("invalid options not caught %v", err)
	}

	options.fetcher = fetcher
	err = process_annotation_types(options)
	if err == nil || err.Error() != "no annotation types to process" {
		t.Error("process did not return error with bad type array")
	}

	options.annotation_types = annotation_types

	annotation := Annotation{
		Object_id:       "smevent:campaignID:eventID",
		Reference_id:    "qcr.app.dev",
		Annotation_type: "name",
		Value:           "an event name",
		Annotator:       "alex"}

	fetcher.annotations = []Annotation{annotation}

	err = process_annotation_types(options)
	if err != nil {
		t.Errorf("no error expected, process failed with error: %v", err)
	}

}

func Test_process_annotations(t *testing.T) {
	//func process_annotations(annotations []Annotation, pagerFactory LoogoPagerFactory) error {

	annotations := []Annotation{{
		Object_id:       "smevent:campaignID:eventID",
		Reference_id:    "qcr.app.dev",
		Annotation_type: "name",
		Value:           "an event name",
		Annotator:       "alex"}}

	pagerFactory := mock_pager_factory{}

	err := process_annotations(annotations, pagerFactory)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

type mockFetcher struct {
	annotations []Annotation
}

func (af mockFetcher) fetch(options annotationOptions) ([]Annotation, error) {
	return af.annotations, nil
}

type mock_pager_factory struct{}

func (pf mock_pager_factory) generate(params loogo.NewPagerParams) (loogo.PagerInterface, error) {
	return mock_pager{}, nil
}

type mock_pager struct{}

func (p mock_pager) GetNext() (loogo.Docs, error) {
	return loogo.Docs{
		loogo.Doc{
			"campaign": "string",
			"event":    "string",
			"features": []string{"string"},
			"name":     "string",
			"relevant": true,
			"id":       "string",
		},
	}, nil
}

func (p mock_pager) PageOver(docFunc func(doc loogo.Doc, done func())) error {
	return nil
}
