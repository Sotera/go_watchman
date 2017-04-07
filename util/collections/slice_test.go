package collections

import (
	"reflect"
	"testing"
)

func TestPrepend(t *testing.T) {
	t.Run("prepend", func(t *testing.T) {
		s := []interface{}{2, 1}
		p := Prepend(s, 4, 3)
		got, want := reflect.DeepEqual(p, []interface{}{4, 3, 2, 1}), true
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
