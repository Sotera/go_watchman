package mockery

import (
	"testing"
)

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

func TestGenerateStr(t *testing.T) {
	type args struct {
		n      int
		length int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"smoke test",
			args{
				2, 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(generateStr(tt.args.n, tt.args.length)); got != tt.args.n {
				t.Errorf("generateStr() = %v, want %v", got, tt.args.n)
			}
			if got := len(generateStr(tt.args.n, tt.args.length)[0]); got != tt.args.length {
				t.Errorf("generateStr() = %v, want %v", got, tt.args.length)
			}
		})
	}
}
