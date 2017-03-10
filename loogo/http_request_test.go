package loogo

import "testing"

type TestHTTPClient struct {
}

func (c *TestHTTPClient) DoRequest(params NewRequestParams) ([]byte, error) {
	return []byte(`{"name":"riot","things":["a","b"]}`), nil
}

func TestNewRequest(t *testing.T) {
	doc := Doc{}
	url := "http://bogus"
	r := HTTPRequestParser{
		client: &TestHTTPClient{},
	}
	err := r.NewRequest(NewRequestParams{URL: url}, &doc)
	if err != nil {
		t.Error(err)
	}
	got := doc["name"]
	want := "riot"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestNewRequest_custom_struct(t *testing.T) {
	type MyObj struct {
		Name   string   `json:"name"`
		Things []string `json:"things"`
	}
	doc := MyObj{}
	url := "http://bogus"
	r := HTTPRequestParser{
		client: &TestHTTPClient{},
	}
	err := r.NewRequest(NewRequestParams{URL: url}, &doc)
	if err != nil {
		t.Error(err)
	}
	got := doc.Name
	want := "riot"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
