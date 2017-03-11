package loogo

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
)

// PagerInterface helps with dependency injection-based testing. see example dir.
type PagerInterface interface {
	PageOver(func(Doc, func())) error
	GetNext() (Docs, error)
}

// Docs are docs returned from query.
type Docs []Doc

// Doc is a single item.
type Doc map[string]interface{}

// CountDoc returns from a /count endpoint.
type CountDoc struct {
	Count int `json:"count"`
}

// QueryParams is a slice of QueryParam.
type QueryParams []QueryParam

// QueryParam is a field-value pair.
type QueryParam struct {
	QueryType string
	Field     string
	Values    []string
}

// APIErrorDoc matches API nested error object
type APIErrorDoc struct {
	Error APIError `json:"error"`
}

// error docs have differing data structure. search for err message.
func (doc *APIErrorDoc) getMessage() string {
	// should only be one of these
	return doc.Error.Message + doc.Error.ErrMsg
}

// APIError matches API error fields
type APIError struct {
	Name string `json:"name"`
	// Status     int    `json:"status"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	ErrMsg     string `json:"errmsg"`
}

// Pager pages over docs from URL.
// Uses a scrolling technique, not offsets (offsets are slow in mongo).
type Pager struct {
	URL           string
	CurrentPage   int // 0-based
	TotalPages    int
	TotalCount    int
	TotalReturned int
	PageSize      int
	ScrollID      string
	OrderBy       string
	Query         string
	Parser        RequestParser // for mocking convenience
}

// NewPagerParams are params to NewPager.
type NewPagerParams struct {
	URL      string
	Params   QueryParams
	PageSize int
}

// NewPager inits a Pager instance.
// TODO: NewPager is difficult to test: Unable to mock network calls.
func NewPager(params NewPagerParams) (*Pager, error) {
	if params.PageSize == 0 {
		params.PageSize = 100 // default
	}
	URL := strings.TrimRight(params.URL, "/")
	orderBy := "_id"
	var scrollID string

	// 'findOne', 'count' do not use 'filter' prefix
	qs := buildQuery(params.Params, true)

	countEndpoint := URL + "/count" + qs
	tc, err := getCount(countEndpoint)
	if err != nil {
		return nil, err
	}

	// prep for scrolling if count > 0
	if tc > 0 {
		// find first item matching query, sorted
		findOneEndpoint := URL + "/findone" + qs + fmt.Sprintf("&filter[order]=%s", orderBy)
		scrollID, err = getFirstID(findOneEndpoint)
		if err != nil {
			return nil, err
		}
	}

	return &Pager{
		URL:           URL,
		TotalCount:    tc,
		TotalPages:    tc/params.PageSize + 1,
		PageSize:      params.PageSize,
		OrderBy:       orderBy,
		ScrollID:      scrollID,
		CurrentPage:   0,
		TotalReturned: 0,
		Query:         buildQuery(params.Params, false),
		Parser:        &HTTPRequestParser{Client: &HTTPClient{}},
	}, nil
}

func getCount(countEndpoint string) (int, error) {
	doc := CountDoc{}
	parser := HTTPRequestParser{Client: &HTTPClient{}}
	// this can hide an API error since decoding to CountDoc
	// will be count == 0. maybe thats ok here?
	err := parser.NewRequest(NewRequestParams{URL: countEndpoint}, &doc)
	if err != nil {
		return -1, err
	}

	return int(doc.Count), nil
}

// first ID in set is needed to seed the scrolling.
func getFirstID(findOneEndpoint string) (string, error) {
	doc := Doc{}
	client := &HTTPClient{}

	// store body for multiple use
	body, err := client.DoRequest(NewRequestParams{URL: findOneEndpoint})
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &doc)
	if err != nil {
		return "", err
	}

	//check for error message from API
	_, ok := doc["error"]
	if ok {
		errRes := APIErrorDoc{}
		err = json.Unmarshal(body, &errRes)
		if err != nil {
			return "", err
		}

		return "", errors.New(errRes.Error.Message)
	}

	return doc["id"].(string), nil
}

// byPage adds scrolling filters.
func (p *Pager) byPage() string {
	op := "gte"
	if p.CurrentPage > 0 {
		op = "gt"
	}

	return fmt.Sprintf("filter[where][id][%s]=%s&filter[limit]=%d&filter[order]=%s", op, p.ScrollID, p.PageSize, p.OrderBy)
}

// GetNext returns next set of docs, nil if past upper bounds, or any error.
func (p *Pager) GetNext() (Docs, error) {
	docs := Docs{}

	if p.TotalCount == 0 || (p.TotalCount == p.TotalReturned) {
		return docs, nil
	}

	// always update pager state
	defer func() {
		p.CurrentPage++

		// set last item in batch as new scrollid
		if len(docs) > 0 {
			p.ScrollID = docs[len(docs)-1]["id"].(string)
			p.TotalReturned += len(docs)
		}
	}()

	url := strings.Join([]string{p.URL + p.Query, p.byPage()}, "&")

	err := p.Parser.NewRequest(NewRequestParams{URL: url}, &docs)
	if err != nil {
		return nil, err
	}

	return docs, nil
}

// PageOver runs docFunc for each doc in each page.
// You must call done() in your docFunc, at the end of the func body,
// to indicate async goroutine as complete.
func (p *Pager) PageOver(docFunc func(doc Doc, done func())) error {
	var wg sync.WaitGroup
	wg.Add(p.TotalCount)

	done := func() { wg.Done() }

	for {
		docs, err := p.GetNext()
		if err != nil {
			return err
		}
		if len(docs) == 0 {
			break
		}

		for _, doc := range docs {
			go docFunc(doc, done)
			// wg.Done()
		}
	}

	wg.Wait()
	return nil
}

// TestPager is a bare bones pager with bogus docs.
type TestPager struct{}

// GetNext for tests
func (p *TestPager) GetNext() (Docs, error) {
	return Docs{
		Doc{
			"id":   "123",
			"name": "apple",
		},
		Doc{
			"id":   "456",
			"name": "pear",
		},
	}, nil
}

// PageOver for tests
func (p *TestPager) PageOver(docFunc func(doc Doc, done func())) error {
	return nil
}
