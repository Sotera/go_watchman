package loogo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// HTTPRequester helps with dependency injection-based testing.
type HTTPRequester interface {
	NewRequest(NewRequestParams, interface{}) error
}

// NewRequestParams are params to NewRequest.
// Accepts same methods as http std lib.
type NewRequestParams struct {
	URL        string
	Params     QueryParams
	HTTPMethod string
	Body       []byte
}

// DoRequest wraps http client Do() to handle any type of http method.
// Defaults to GET request.
func DoRequest(params NewRequestParams) ([]byte, error) {
	params.URL = strings.TrimRight(params.URL, "/")
	if params.HTTPMethod == "" {
		params.HTTPMethod = "GET"
	}

	client := &http.Client{}

	req, err := http.NewRequest(
		params.HTTPMethod,
		params.URL+BuildQuery(params.Params, false),
		bytes.NewBuffer(params.Body),
	)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

// NewRequest populates result arg with returned values.
// For GET requests, result should be a Docs{}. Otherwise, result should
// be a Doc{} since PUT,POST,DELETE requests return a single item.
// The result arg allows the client to tell this function what to expect
// in an effort to reduce code complexity here.
// Preferably the client sends a struct so that it can easily distinguish
// between a valid result and an api error message.
func NewRequest(params NewRequestParams, result interface{}) error {
	body, err := DoRequest(params)
	if err != nil {
		return err
	}

	// TODO: this will not distinguish b/n valid results and err doc if
	// result is generic. do we care about that here?
	err = json.Unmarshal(body, result)
	if err != nil {
		return errors.New(string(body))
	}

	return nil
}
