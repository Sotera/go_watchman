package loogo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// QueryParams is a slice of QueryParam
type QueryParams []QueryParam

// QueryParam is a field-value pair
type QueryParam struct {
	QueryType string
	Field     string
	Values    []string
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
		default:
			fmt.Println("unknown QueryType")
		}
	}

	return "?" + strings.Join(qs, "&")
}

// Between returns loopback querystring for between queries
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

// Eq returns loopback querystring for equality queries
func Eq(p QueryParam, countOnly bool) string {
	prefix := "filter"
	if countOnly {
		prefix = ""
	}

	return fmt.Sprintf("%s[where][%s]=%s", prefix, p.Field, p.Values[0])
}

// Pager pages over results from URL.
// Uses a scrolling technique, not offsets (offsets are slow in mongo).
type Pager struct {
	URL         string
	CurrentPage int
	TotalPages  int
	TotalCount  int
	PageSize    int
	ScrollID    string
	OrderBy     string
}

// NewPager inits a Pager instance
func NewPager(URL string, params QueryParams, pageSize int) *Pager {
	URL = strings.TrimRight(URL, "/")
	orderBy := "_id ASC"

	countEndpoint := URL + BuildQuery(params, true)
	tc := getCount(countEndpoint)

	findOneEndpoint := URL + "/findone" + BuildQuery(params, false) + fmt.Sprintf("&filter[order]=%s", orderBy)

	scrollID := getFirstID(findOneEndpoint)

	return &Pager{
		URL:        URL,
		TotalCount: tc,
		TotalPages: tc/pageSize + 1,
		PageSize:   pageSize,
		OrderBy:    orderBy,
		ScrollID:   scrollID,
	}
}

func getCount(countEndpoint string) int {
	resp, err := http.Get(countEndpoint)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	results := map[string]float64{}

	json.NewDecoder(resp.Body).Decode(&results)

	return int(results["count"])
}

func getFirstID(findOneEndpoint string) string {
	resp, err := http.Get(findOneEndpoint)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	result := Result{}

	json.NewDecoder(resp.Body).Decode(&result)

	return result["id"].(string)
}

// ByPage adds scrolling filters
func (p *Pager) ByPage() string {
	op := "gte"
	if p.CurrentPage > 0 {
		op = "gt"
	}

	return fmt.Sprintf("filter[where][id][%s]=%s&filter[limit]=%d&filter[order]=%s", op, p.ScrollID, p.PageSize, p.OrderBy)
}

// Results are query results
type Results []Result

// Result is a single item
type Result map[string]interface{}

// GetNext returns next page of results, nil if past upper bounds
func (p *Pager) GetNext() Results {
	if p.TotalCount == 0 || (p.CurrentPage >= p.TotalPages) {
		return nil
	}

	url := strings.Join([]string{p.URL, p.ByPage()}, "&")

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	results := Results{}

	json.NewDecoder(resp.Body).Decode(&results)
	fmt.Println(len(results))

	if len(results) != 0 {
		// get last item in batch for cursor marker
		p.ScrollID = results[len(results)-1]["id"].(string)
	} else {
		return nil
	}

	p.CurrentPage++
	return results
}
