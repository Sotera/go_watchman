package follow_along

import (
	"reflect"
	"testing"
)

func TestSet_Items(t *testing.T) {
	type fields struct {
		items map[interface{}]bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []interface{}
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Set{
				items: tt.fields.items,
			}
			if got := s.Items(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set.Items() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_add(t *testing.T) {
	type fields struct {
		items map[interface{}]bool
	}
	type args struct {
		item interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Set{
				items: tt.fields.items,
			}
			s.add(tt.args.item)
		})
	}
}