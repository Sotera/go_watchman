package go_watchman

import "testing"

type mockFetcher struct {
	annotations []Annotation
}

func (af mockFetcher) fetch(options annotationOptions) (*[]Annotation, error) {
	return &af.annotations, nil
}

func TestFetchAnnotations(t *testing.T) {
	options := annotationOptions{
		startTime:      "",
		endTime:        "",
		apiRoot:        "",
		annotationType: ""}
	fetcher := mockFetcher{}

	annotation := Annotation{
		Object_id:       "smevent:campaignID:eventID",
		Reference_id:    "qcr.app.dev",
		Annotation_type: "name",
		Value:           "an event name",
		Annotator:       "alex"}

	fetcher.annotations = []Annotation{annotation}

	val, err := fetch_annotations(options, fetcher)

	if err != nil {
		t.Errorf("fetch failed with error: %v", err)
	}

	if len(*val) != 1 {
		t.Errorf("annotation count expectd %v got %v", 0, len(*val))
	}
}

func TestProcessAnnotationTypes(t *testing.T) {

	fetcher := mockFetcher{}

	annotation_types := []string{"name", "relevance"}

	options := annotationOptions{
		startTime:      "test",
		endTime:        "test",
		apiRoot:        "test",
		annotationType: "test"}

	err := process_annotation_types(options, annotation_types, nil)
	if err == nil || err.Error() != "fetcher instance was nil, please provide a fetcher" {
		t.Errorf("fetcher instance was nil not caught: %v", err)
	}

	err = process_annotation_types(annotationOptions{}, annotation_types, fetcher)
	if err == nil || err.Error() != "invalid options" {
		t.Errorf("invalid options not caught %v", err)
	}

	err = process_annotation_types(options, nil, fetcher)
	if err == nil || err.Error() != "no annotation types to process" {
		t.Error("process did not return error with bad type array")
	}

	err = process_annotation_types(options, annotation_types, fetcher)
	if err != nil {
		t.Errorf("no error expected, process failed with error: %v", err)
	}

}
