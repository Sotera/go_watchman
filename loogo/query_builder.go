package loogo

import (
	"fmt"
	"strings"
)

// buildQuery returns combined query string from QueryParams.
// Prepends '?' to return value.
// countOnly: don't include 'filter' in query string.
func buildQuery(params QueryParams, countOnly bool) string {
	qs := []string{}

	for _, p := range params {
		switch strings.ToLower(p.QueryType) {
		case "between":
			qs = append(qs, between(p, countOnly))
		case "eq":
			qs = append(qs, eq(p, countOnly))
		case "inq":
			qs = append(qs, inq(p, countOnly))
		default:
			fmt.Println("unknown QueryType")
		}
	}

	if len(qs) == 0 {
		return ""
	}
	return "?" + strings.Join(qs, "&")
}

// between returns loopback querystring for between queries.
func between(p QueryParam, countOnly bool) string {
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

// eq returns loopback querystring for equality queries.
func eq(p QueryParam, countOnly bool) string {
	prefix := "filter"
	if countOnly {
		prefix = ""
	}

	return fmt.Sprintf("%s[where][%s]=%s", prefix, p.Field, p.Values[0])
}

// inq returns loopback querystring for inclusion queries.
// ex. filter[where][name][inq]=foo&filter[where][name][inq]=bar
func inq(p QueryParam, countOnly bool) string {
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
