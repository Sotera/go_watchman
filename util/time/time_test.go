package time

import "testing"
import "time"

func TestStrToUnixMs(t *testing.T) {
	type args struct {
		layout  string
		timeStr string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			"happy",
			args{
				layout:  time.RFC3339,
				timeStr: "2017-01-01T00:00:00Z",
			},
			1483228800000,
			false,
		},
		{
			"sad",
			args{
				layout:  time.RFC3339,
				timeStr: "bogus",
			},
			-1,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StrToUnixMs(tt.args.layout, tt.args.timeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrToUnixMs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StrToUnixMs() = %v, want %v", got, tt.want)
			}
		})
	}
}
