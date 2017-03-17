package annotations

import (
	"strings"
	"testing"

	"github.com/Sotera/go_watchman/loogo"
)

func TestFetchAnnotations(t *testing.T) {
	fetcher := MockFetcher{}

	pagerFactory := MockPagerFactory{}

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
	pagerFactory := MockPagerFactory{}
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

func TestProcessAnnotations(t *testing.T) {
	//func process_annotations(annotations []Annotation, pagerFactory LoogoPagerFactory) error {
	fetcher := MockFetcher{}

	testAnnos := []Annotation{{
		ObjectID:       "smevent:campaignID:eventID",
		ReferenceID:    "qcr.app.dev",
		AnnotationType: "label",
		Value:          "an event name",
		Annotator:      "alex"}}

	pagerFactory := MockPagerFactory{}
	parserFactory := MockParserFactory{}
	options := AnnotationOptions{
		StartTime:         "",
		EndTime:           "",
		AnnotationAPIRoot: "",
		AnnotationType:    "",
		AnnotationTypes:   []string{"label", "relevant"},
		Fetcher:           fetcher,
	}

	err := ProcessAnnotations(testAnnos, options)
	if err == nil {
		t.Errorf("error expected: %v", err)
	}

	options.PagerFactory = pagerFactory

	err = ProcessAnnotations(testAnnos, options)
	if err == nil {
		t.Errorf("error expected: %v", err)
	}

	options.ParserFactory = parserFactory
	err = ProcessAnnotations(testAnnos, options)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	testAnnos[0].AnnotationType = "relevant"
	testAnnos[0].Value = "false"
	options.AnnotationType = "relevant"
	options.PagerFactory = MockPagerFactory{
		ReturnEmpty: true,
	}
	err = ProcessAnnotations(testAnnos, options)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

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

type MockPagerFactory struct {
	ReturnEmpty bool
}

func (pf MockPagerFactory) Generate(params loogo.NewPagerParams) (loogo.PagerInterface, error) {
	pager := MockPager{}
	pager.ReturnEmpty = pf.ReturnEmpty
	if strings.Contains(params.URL, "/events") {
		println("sending event")
		pager.ReturnEmpty = false
		pager.ReturnEvent = true
	}

	return pager, nil
}

type MockPager struct {
	ReturnEmpty bool
	ReturnEvent bool
}

func (p MockPager) GetNext() (loogo.Docs, error) {
	if p.ReturnEmpty {
		println("returning empty doc")
		return loogo.Docs{}, nil
	}

	if p.ReturnEvent {
		println("returning event doc")
		return loogo.Docs{
			loogo.Doc{
				"id": "eventID",
			},
		}, nil
	}

	println("returning annotation model doc")
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

type MockParserFactory struct{}

func (pf MockParserFactory) Generate() loogo.RequestParser {
	return &MockParser{}
}

type MockParser struct{}

func (r *MockParser) NewRequest(params loogo.NewRequestParams, result interface{}) error {
	doc := loogo.Doc{}
	doc["_id"] = "eventID"
	doc["hashtags"] = []string{}
	result = doc
	return nil
}

type MockRequester struct{}
