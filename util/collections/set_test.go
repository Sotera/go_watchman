package collections

import (
	"testing"
)

func TestSet_Items(t *testing.T) {
	var got, want interface{}
	t.Run("add", func(t *testing.T) {
		s := &Set{}
		s.Add("a", "b")
		got, want = len(s.Items()), 2
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestSet_Add(t *testing.T) {
	var got, want interface{}
	t.Run("add", func(t *testing.T) {
		s := &Set{}
		s.Add("item")
		s.Add("item", "another")
		got, want = len(s.Items()), 2
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		got, want = s.Items()[0], "item"
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		got, want = s.Items()[1], "another"
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestSet_Delete(t *testing.T) {
	var got, want interface{}
	t.Run("delete", func(t *testing.T) {
		s := &Set{}
		s.Add("a", "b")
		s.Delete("b")
		got, want = len(s.Items()), 1
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		got, want = s.Items()[0], "a"
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
