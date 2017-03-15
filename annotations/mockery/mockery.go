package mockery

import (
	"strconv"
	"strings"
	"time"

	"errors"

	ann "github.com/Sotera/go_watchman/annotations"
	"github.com/Sotera/go_watchman/loogo"
	timeu "github.com/Sotera/go_watchman/util/time"
)

type watchmanEvent map[string]interface{}
type watchmanEvents []watchmanEvent

// Mockery is a fake annotator with a parser dependency.
type Mockery struct {
	Parser loogo.RequestParser
}

// GetAnnotations returns mocked annotations based on existing watchman events.
func (m Mockery) GetAnnotations(options ann.AnnotationOptions) ([]ann.Annotation, error) {
	var err error
	events := watchmanEvents{}
	annotations := []ann.Annotation{}

	startTime, err := timeu.StrToUnixMs(time.RFC3339, options.StartTime)
	if err != nil {
		return nil, err
	}
	endTime, err := timeu.StrToUnixMs(time.RFC3339, options.EndTime)
	if err != nil {
		return nil, err
	}

	err = m.Parser.NewRequest(
		loogo.NewRequestParams{
			URL: options.APIRoot,
			Params: loogo.QueryParams{
				loogo.QueryParam{
					Field:     "start_time_ms",
					QueryType: "between",
					Values: []string{
						strconv.Itoa(startTime),
						strconv.Itoa(endTime),
					},
				},
			},
		},
		&events,
	)
	if err != nil {
		return nil, err
	}

	for _, evt := range events {
		// real campaign data is a list of objects. lets just use a fake list.
		campids := generateStr(3, 5)

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
