package loogo

import "testing"

func Test_Between(t *testing.T) {
	tests := []struct {
		name  string
		param QueryParam
		want  string
	}{
		{
			"between",
			QueryParam{
				QueryType: "Between",
				Field:     "created",
				Values:    []string{"0", "1"},
			},
			"filter[where][created][between][0]=0&filter[where][created][between][1]=1",
		},
		{
			"between",
			QueryParam{
				QueryType: "Between",
				Field:     "created",
				Values:    []string{"0", "1"},
			},
			"[where][created][between][0]=0&[where][created][between][1]=1",
		},
	}
	for _, tt := range tests {
		if got := Between(tt.param, false); got != tt.want {
			t.Errorf("%q. Between() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_Eq(t *testing.T) {
	tests := []struct {
		name  string
		param QueryParam
		want  string
	}{
		{
			"eq",
			QueryParam{
				QueryType: "Eq",
				Field:     "featurizer",
				Values:    []string{"image"},
			},
			"filter[where][featurizer]=image",
		},
		{
			"eq",
			QueryParam{
				QueryType: "Eq",
				Field:     "featurizer",
				Values:    []string{"image"},
			},
			"[where][featurizer]=image",
		},
	}
	for _, tt := range tests {
		if got := Eq(tt.param, false); got != tt.want {
			t.Errorf("%q. Eq() = %v, want %v", tt.name, got, tt.want)
		}
	}

}

func Test_BuildQuery(t *testing.T) {
	tests := []struct {
		name   string
		params []QueryParam
		want   string
	}{
		{
			"buildquery",
			[]QueryParam{
				{
					QueryType: "Eq",
					Field:     "featurizer",
					Values:    []string{"image"},
				},
				{
					QueryType: "Eq",
					Field:     "status",
					Values:    []string{"new"},
				},
			},
			"?filter[where][featurizer]=image&filter[where][status]=new",
		},
	}
	for _, tt := range tests {
		if got := BuildQuery(tt.params, false); got != tt.want {
			t.Errorf("%q. BuildQuery() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
