package annotations

import (
	"github.com/sotera/go_watchman/loogo"
	"testing"
)

func TestFetchAnnotations(t *testing.T) {
	fetcher := MockFetcher{}

	pagerFactory := mockPagerFactory{}

	annotation := Annotation{
		Object_id:       "smevent:campaignID:eventID",
		Reference_id:    "qcr.app.dev",
		Annotation_type: "name",
		Value:           "an event name",
		Annotator:       "alex"}

	fetcher.Annotations = []Annotation{annotation}

	options := AnnotationOptions{
		StartTime:         "",
		EndTime:           "",
		AnnotationApiRoot: "",
		AnnotationType:    "",
		Annotation_types:  []string{"test"},
		Fetcher:           fetcher,
		PagerFactory:      pagerFactory,
	}

	val, err := FetchAnnotations(options)

	if err != nil {
		t.Errorf("fetch failed with error: %v", err)
	}

	if len(val) != 1 {
		t.Errorf("annotation count expectd %v got %v", 1, len(val))
	}
}

func TestProcessAnnotationTypes(t *testing.T) {

	fetcher := MockFetcher{}

	annotation_types := []string{"name", "relevance"}

	options := AnnotationOptions{
		StartTime:         "test",
		EndTime:           "test",
		AnnotationApiRoot: "test",
		AnnotationType:    "test",
		Fetcher:           nil,
		Annotation_types:  nil}

	err := ProcessAnnotationTypes(options)
	if err == nil || err.Error() != "fetcher instance was nil, please provide a fetcher" {
		t.Errorf("fetcher instance was nil not caught: %v", err)
	}

	err = ProcessAnnotationTypes(AnnotationOptions{})
	if err == nil || err.Error() != "invalid options" {
		t.Errorf("invalid options not caught %v", err)
	}

	options.Fetcher = fetcher
	err = ProcessAnnotationTypes(options)
	if err == nil || err.Error() != "no annotation types to process" {
		t.Error("process did not return error with bad type array")
	}

	options.Annotation_types = annotation_types

	annotation := Annotation{
		Object_id:       "smevent:campaignID:eventID",
		Reference_id:    "qcr.app.dev",
		Annotation_type: "name",
		Value:           "an event name",
		Annotator:       "alex"}

	fetcher.Annotations = []Annotation{annotation}

	err = ProcessAnnotationTypes(options)
	if err != nil {
		t.Errorf("no error expected, process failed with error: %v", err)
	}

}

func TestProcessAnnotations(t *testing.T) {
	//func process_annotations(annotations []Annotation, pagerFactory LoogoPagerFactory) error {

	annotations := []Annotation{{
		Object_id:       "smevent:campaignID:eventID",
		Reference_id:    "qcr.app.dev",
		Annotation_type: "name",
		Value:           "an event name",
		Annotator:       "alex"}}

	pagerFactory := mockPagerFactory{}

	err := ProcessAnnotations(annotations, pagerFactory)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseAnnotationId(t *testing.T) {

	campaignId1, eventId1 := ParseAnnotationId("smevent:campaignID:eventID")
	if campaignId1 != "campaignID" || eventId1 != "eventID" {
		t.Errorf("unexpected error: %v %v", campaignId1, eventId1)
	}

	campaignId2, eventId2 := ParseAnnotationId("smevent::eventID")
	if campaignId2 != "" || eventId2 != "eventID" {
		t.Errorf("unexpected error: %v %v", campaignId2, eventId2)
	}

	campaignId3, eventId3 := ParseAnnotationId("smevent::")
	if campaignId3 != "" || eventId3 != "" {
		t.Errorf("unexpected error: %v %v", campaignId3, eventId3)
	}

	campaignId4, eventId4 := ParseAnnotationId("smevent:")
	if campaignId4 != "" || eventId4 != "" {
		t.Errorf("unexpected error: %v %v", campaignId4, eventId4)
	}

	campaignId5, eventId5 := ParseAnnotationId("")
	if campaignId5 != "" || eventId5 != "" {
		t.Errorf("unexpected error: %v %v", campaignId5, eventId5)
	}

	campaignId6, eventId6 := ParseAnnotationId("smevent::eventID")
	if campaignId6 != "" || eventId6 != "eventID" {
		t.Errorf("unexpected error: %v %v", campaignId6, eventId6)
	}
}

type MockFetcher struct {
	Annotations []Annotation
}

func (af MockFetcher) Fetch(options AnnotationOptions) ([]Annotation, error) {
	return af.Annotations, nil
}

type mockPagerFactory struct{}

func (pf mockPagerFactory) generate(params loogo.NewPagerParams) (loogo.PagerInterface, error) {
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
