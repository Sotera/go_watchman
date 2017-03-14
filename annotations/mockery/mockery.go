package mockery

import (
	"strings"

	"errors"

	ann "github.com/Sotera/go_watchman/annotations"
	"github.com/Sotera/go_watchman/loogo"
)

type watchmanEvent map[string]interface{}
type watchmanEvents []watchmanEvent

// Mockery is a fake annotator with a parser dependency.
type Mockery struct {
	Parser loogo.RequestParser
}

// GetAnnotations returns mocked annotations based on existing watchman events.
func (m Mockery) GetAnnotations(options ann.AnnotationOptions) ([]ann.Annotation, error) {

	events := watchmanEvents{}

	annotations := []ann.Annotation{}

	err := m.Parser.NewRequest(
		loogo.NewRequestParams{
			URL: options.APIRoot,
			Params: loogo.QueryParams{
				loogo.QueryParam{
					Field:     "start_time_ms",
					QueryType: "between",
					Values:    []string{options.StartTime, options.EndTime},
				},
			},
		},
		&events,
	)
	if err != nil {
		return nil, err
	}

	for _, evt := range events {
		campids := evt["campaigns"].([]string)
		if len(campids) == 0 {
			campids = []string{"123", "456"}
		}

		for _, cid := range campids {
			var val string
			if options.AnnotationType == "relevant" {
				val = takeOne([]string{"true", "false"})
			} else if options.AnnotationType == "name" {
				val = takeOne([]string{})
			} else {
				return nil, errors.New("unknown annotation type:" + options.AnnotationType)
			}

			annotations = append(annotations,
				ann.Annotation{
					AnnotationType: options.AnnotationType,
					Annotator:      "pepe",
					CampaignID:     cid,
					EventID:        evt["id"].(string),
					ObjectID:       strings.Join([]string{"smevent", evt["id"].(string), cid}, ":"),
					ReferenceID:    "qcr.app.dev",
					Value:          val,
				},
			)
		}
	}

	return annotations, nil
}
