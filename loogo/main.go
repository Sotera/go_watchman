package loogo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

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

// APIError matches API error fields
type APIError struct {
	Name string `json:"name"`
	// Status     int    `json:"status"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// BuildQuery returns combined query string from QueryParams.
// Prepends '?' to return value.
// countOnly: don't include 'filter' in query string.
func BuildQuery(params QueryParams, countOnly bool) string {
	qs := []string{}

	for _, p := range params {
		switch strings.ToLower(p.QueryType) {
		case "between":
			qs = append(qs, Between(p, countOnly))
		case "eq":
			qs = append(qs, Eq(p, countOnly))
		case "inq":
			qs = append(qs, Inq(p, countOnly))
		default:
			fmt.Println("unknown QueryType")
		}
	}

	return "?" + strings.Join(qs, "&")
}

// Between returns loopback querystring for between queries.
func Between(p QueryParam, countOnly bool) string {
	qs := []string{}
	prefix := "filter"
	if countOnly {
		prefix = ""
	}

	for i, v := range p.Values {
		qs = append(qs, fmt.Sprintf("[where][%s][between][%d]=%s", p.Field, i, v))
	}

	parts := []string{prefix, strings.Join(qs, "&"+prefix)}

	return strings.Join(parts, "")
}

// Eq returns loopback querystring for equality queries.
func Eq(p QueryParam, countOnly bool) string {
	prefix := "filter"
	if countOnly {
		prefix = ""
	}

	return fmt.Sprintf("%s[where][%s]=%s", prefix, p.Field, p.Values[0])
}

// Inq returns loopback querystring for inclusion queries.
// ex. filter[where][name][inq]=foo&filter[where][name][inq]=bar
func Inq(p QueryParam, countOnly bool) string {
	qs := []string{}
	prefix := "filter"
	if countOnly {
		prefix = ""
	}

	for _, v := range p.Values {
		qs = append(qs, fmt.Sprintf("[where][%s][inq]=%s", p.Field, v))
	}

	parts := []string{prefix, strings.Join(qs, "&"+prefix)}

	return strings.Join(parts, "")
}

// Pager pages over docs from URL fetch.
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
}

// NewPagerParams are params to NewPager.
type NewPagerParams struct {
	URL      string
	Params   QueryParams
	PageSize int
}

// NewPager inits a Pager instance.
func NewPager(params NewPagerParams) (*Pager, error) {
	if params.PageSize == 0 {
		params.PageSize = 100 // default
	}
	URL := strings.TrimRight(params.URL, "/")
	orderBy := "_id"
	var scrollID string

	// 'findOne', 'count' do not use 'filter' prefix
	qs := BuildQuery(params.Params, true)

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
		Query:         BuildQuery(params.Params, false),
	}, nil
}

func getCount(countEndpoint string) (int, error) {
	resp, err := http.Get(countEndpoint)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	doc := CountDoc{}

	// this can hide an API error since decoding to CountDoc
	// will be count == 0. maybe thats ok here?
	err = json.NewDecoder(resp.Body).Decode(&doc)
	if err != nil {
		return -1, err
	}

	return int(doc.Count), nil
}

func getFirstID(findOneEndpoint string) (string, error) {
	resp, err := http.Get(findOneEndpoint)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	doc := Doc{}

	// store body for multiple use
	body, err := ioutil.ReadAll(resp.Body)
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

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&docs)
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

type PagerInterface interface {
	PageOver(func(Doc, func())) error
	GetNext() (Docs, error)
}

type TestPager struct{}

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

func (p *TestPager) PageOver(docFunc func(doc Doc, done func())) error {
	return nil
}
