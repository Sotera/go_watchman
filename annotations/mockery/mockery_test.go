package mockery

import (
	"reflect"
	"testing"

	a "github.com/Sotera/go_watchman/annotations"
	"github.com/Sotera/go_watchman/loogo"
)

type TestHTTPClient struct {
}

func (c *TestHTTPClient) DoRequest(params loogo.NewRequestParams) ([]byte, error) {
	return []byte(`{"name":"riot","things":["a","b"]}`), nil
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
		want    []a.Annotation
		wantErr bool
	}{
		{
			"smoke test",
			args{
				a.AnnotationOptions{},
			},
			[]a.Annotation{},
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAnnotations() = %v, want %v", got, tt.want)
			}
		})
	}
}
