package mockery

import (
	"testing"

	a "github.com/Sotera/go_watchman/annotations"
	"github.com/Sotera/go_watchman/loogo"
)

type TestHTTPClient struct {
}

func (c *TestHTTPClient) DoRequest(params loogo.NewRequestParams) ([]byte, error) {
	return []byte(`[{"id":"123"}]`), nil
}

func TestGetAnnotations(t *testing.T) {
	parser := &loogo.HTTPRequestParser{
		Client: &TestHTTPClient{},
	}
	mocker := Mockery{
		Parser: parser,
	}

	type args struct {
		options a.AnnotationOptions
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"smoke test",
			args{
				a.AnnotationOptions{
					StartTime:      "2017-03-01T00:00:00Z",
					EndTime:        "2017-03-02T00:00:00Z",
					AnnotationType: "name",
				},
			},
			"123",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mocker.GetAnnotations(tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAnnotations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got[0].EventID != tt.want {
				t.Errorf("GetAnnotations() = %v, want %v", got[0].EventID, tt.want)
			}
		})
	}
}
