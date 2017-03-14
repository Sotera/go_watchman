package mockery

import "testing"

func TestTakeOne(t *testing.T) {
	type args struct {
		vals []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"smoke test",
			args{
				[]string{"true", "true"},
			},
			"true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := takeOne(tt.args.vals); got != tt.want {
				t.Errorf("takeOne() = %v, want %v", got, tt.want)
			}
		})
	}
}
