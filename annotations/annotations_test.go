package annotations

import (
	"testing"

	"github.com/Sotera/go_watchman/loogo"
)

func TestFetchAnnotations(t *testing.T) {
	fetcher := MockFetcher{}

	annotation := Annotation{
		ObjectID:       "smevent:campaignID:eventID",
		ReferenceID:    "qcr.app.dev",
		AnnotationType: "label",
		Value:          "an event name",
		Annotator:      "alex"}

	fetcher.Annotations = []Annotation{annotation}

	options := AnnotationOptions{
		APIRoot:           "http://test.com",
		StartTime:         "",
		EndTime:           "",
		AnnotationAPIRoot: "",
		AnnotationType:    "",
		AnnotationTypes:   []string{"test"},
		Fetcher:           fetcher,
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

	annotationTypes := []string{"label", "relevance"}

	options := AnnotationOptions{
		APIRoot:           "http://test.com",
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

	options.AnnotationTypes = annotationTypes

	annotation := Annotation{
		ObjectID:       "smevent:campaignID:eventID",
		ReferenceID:    "qcr.app.dev",
		AnnotationType: "label",
		Value:          "an event name",
		Annotator:      "alex",
		EventID:        "eventID",
		CampaignID:     "campaignID",
	}

	fetcher.Annotations = []Annotation{annotation}

	err = ProcessAnnotationTypes(options)
	if err != nil {
		t.Errorf("no error expected, process failed with error: %v", err)
	}

}

func TestLoogoInterfaces(t *testing.T) {
	pager := loogo.TestPager{}
	page, err := pager.GetNext()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if page[0]["id"].(string) != "123" {
		t.Errorf("incorrect value")
	}
}

func TestProcessAnnotations(t *testing.T) {
	fetcher := MockFetcher{}

	testAnnos := []Annotation{{
		ObjectID:       "smevent:campaignID:eventID",
		ReferenceID:    "qcr.app.dev",
		AnnotationType: "label",
		Value:          "an event name",
		Annotator:      "alex"}}

	options := AnnotationOptions{
		APIRoot:           "http://test.com",
		StartTime:         "",
		EndTime:           "",
		AnnotationAPIRoot: "",
		AnnotationType:    "",
		AnnotationTypes:   []string{"label", "relevant"},
		Fetcher:           fetcher,
	}

	createPager := func(params loogo.NewPagerParams) (loogo.PagerInterface, error) {
		return &MockPager{}, nil
	}
	parser := &MockParser{}
	am := &AnnotationMaker{createPager, parser}

	err := am.ProcessAnnotations(testAnnos, options)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	testAnnos[0].AnnotationType = "relevant"
	testAnnos[0].Value = "false"
	options.AnnotationType = "relevant"

	err = am.ProcessAnnotations(testAnnos, options)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseAnnotationId(t *testing.T) {

	campaignID1, eventID1 := ParseAnnotationID("smevent:campaignID:eventID")
	if campaignID1 != "campaignID" || eventID1 != "eventID" {
		t.Errorf("unexpected error: %v %v", campaignID1, eventID1)
	}

	campaignID2, eventID2 := ParseAnnotationID("smevent::eventID")
	if campaignID2 != "" || eventID2 != "eventID" {
		t.Errorf("unexpected error: %v %v", campaignID2, eventID2)
	}

	campaignID3, eventID3 := ParseAnnotationID("smevent::")
	if campaignID3 != "" || eventID3 != "" {
		t.Errorf("unexpected error: %v %v", campaignID3, eventID3)
	}

	campaignID4, eventID4 := ParseAnnotationID("smevent:")
	if campaignID4 != "" || eventID4 != "" {
		t.Errorf("unexpected error: %v %v", campaignID4, eventID4)
	}

	campaignID5, eventID5 := ParseAnnotationID("")
	if campaignID5 != "" || eventID5 != "" {
		t.Errorf("unexpected error: %v %v", campaignID5, eventID5)
	}

	campaignID6, eventID6 := ParseAnnotationID("smevent::eventID")
	if campaignID6 != "" || eventID6 != "eventID" {
		t.Errorf("unexpected error: %v %v", campaignID6, eventID6)
	}
}

type MockFetcher struct {
	Annotations []Annotation
}

func (af MockFetcher) Fetch(options AnnotationOptions) ([]Annotation, error) {
	return af.Annotations, nil
}

type MockPager struct {
	ReturnEmpty bool
	ReturnEvent bool
}

func (p MockPager) GetNext() (loogo.Docs, error) {
	if p.ReturnEmpty {
		// println("returning empty doc")
		return loogo.Docs{}, nil
	}

	if p.ReturnEvent {
		// println("returning event doc")
		return loogo.Docs{
			loogo.Doc{
				"id": "eventID",
			},
		}, nil
	}

	// println("returning annotation model doc")
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

func (p MockPager) PageOver(docFunc func(doc loogo.Doc, done func())) error {
	return nil
}

type MockParser struct{}

func (r *MockParser) NewRequest(params loogo.NewRequestParams, result interface{}) error {
	doc := loogo.Doc{}
	doc["_id"] = "eventID"
	doc["hashtags"] = [][]interface{}{}
	result = doc
	return nil
}
