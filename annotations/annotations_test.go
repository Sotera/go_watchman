package annotations

import (
	"testing"

	"github.com/Sotera/go_watchman/loogo"
)

func TestFetchAnnotations(t *testing.T) {
	fetcher := MockFetcher{}

	pagerFactory := mockPagerFactory{}

	annotation := Annotation{
		ObjectID:       "smevent:campaignID:eventID",
		ReferenceID:    "qcr.app.dev",
		AnnotationType: "label",
		Value:          "an event name",
		Annotator:      "alex"}

	fetcher.Annotations = []Annotation{annotation}

	options := AnnotationOptions{
		StartTime:         "",
		EndTime:           "",
		AnnotationAPIRoot: "",
		AnnotationType:    "",
		AnnotationTypes:   []string{"test"},
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

	annotation_types := []string{"label", "relevance"}

	options := AnnotationOptions{
		StartTime:         "test",
		EndTime:           "test",
		AnnotationAPIRoot: "test",
		AnnotationType:    "test",
		Fetcher:           nil,
		AnnotationTypes:   nil}

	err := ProcessAnnotationTypes(options)
	if err == nil || err.Error() != "fetcher instance was nil, please provide a fetcher" {
		t.Errorf("fetcher instance was nil not caught: %v", err)
	}

	options.Fetcher = fetcher
	err = ProcessAnnotationTypes(options)
	if err == nil || err.Error() != "no annotation types to process" {
		t.Error("process did not return error with bad type array")
	}

	options.AnnotationTypes = annotation_types

	annotation := Annotation{
		ObjectID:       "smevent:campaignID:eventID",
		ReferenceID:    "qcr.app.dev",
		AnnotationType: "label",
		Value:          "an event name",
		Annotator:      "alex"}

	fetcher.Annotations = []Annotation{annotation}

	err = ProcessAnnotationTypes(options)
	if err != nil {
		t.Errorf("no error expected, process failed with error: %v", err)
	}

}

func TestLoogoInterfaces(t *testing.T) {
	pagerFactory := mockPagerFactory{}
	pager, err := pagerFactory.Generate(loogo.NewPagerParams{
		URL:      "http://localhost:3003/api/annotations",
		Params:   nil,
		PageSize: 1,
	})
	page, err := pager.GetNext()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	value := page[0]["id"].(string)
	println(value)
}

/*func TestProcessAnnotations(t *testing.T) {
	//func process_annotations(annotations []Annotation, pagerFactory LoogoPagerFactory) error {

	testAnnos := []Annotation{{
		Object_id:       "smevent:campaignID:eventID",
		Reference_id:    "qcr.app.dev",
		Annotation_type: "label",
		Value:           "an event name",
		Annotator:       "alex"}}

	pagerFactory := mockPagerFactory{}
	options := AnnotationOptions{
		StartTime:         "",
		EndTime:           "",
		AnnotationApiRoot: "",
		AnnotationType:    "",
		Annotation_types:  []string{"test"},
		Fetcher:           nil,
		PagerFactory:      pagerFactory,
	}

	err := ProcessAnnotations(testAnnos, options)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
*/
func TestParseAnnotationId(t *testing.T) {

	campaignId1, eventId1 := ParseAnnotationID("smevent:campaignID:eventID")
	if campaignId1 != "campaignID" || eventId1 != "eventID" {
		t.Errorf("unexpected error: %v %v", campaignId1, eventId1)
	}

	campaignId2, eventId2 := ParseAnnotationID("smevent::eventID")
	if campaignId2 != "" || eventId2 != "eventID" {
		t.Errorf("unexpected error: %v %v", campaignId2, eventId2)
	}

	campaignId3, eventId3 := ParseAnnotationID("smevent::")
	if campaignId3 != "" || eventId3 != "" {
		t.Errorf("unexpected error: %v %v", campaignId3, eventId3)
	}

	campaignId4, eventId4 := ParseAnnotationID("smevent:")
	if campaignId4 != "" || eventId4 != "" {
		t.Errorf("unexpected error: %v %v", campaignId4, eventId4)
	}

	campaignId5, eventId5 := ParseAnnotationID("")
	if campaignId5 != "" || eventId5 != "" {
		t.Errorf("unexpected error: %v %v", campaignId5, eventId5)
	}

	campaignId6, eventId6 := ParseAnnotationID("smevent::eventID")
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

func (pf mockPagerFactory) Generate(params loogo.NewPagerParams) (loogo.PagerInterface, error) {
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
