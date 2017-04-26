package main

import (
	"reflect"
	"testing"
)

func Test_cache_hit(t *testing.T) {
	c := cache{maxSize: 1, items: []cacheItem{{"three", 3}}}
	got, _ := c.item("three")
	want := 3
	if got.value != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func Test_cache_miss(t *testing.T) {
	c := cache{maxSize: 1, items: []cacheItem{{"three", 3}}}
	_, got := c.item("bogus")
	if got != false {
		t.Errorf("got %v, want %v", got, false)
	}
}

func Test_cache_add(t *testing.T) {
	c := cache{maxSize: 1, items: []cacheItem{{"three", 3}, {"four", 4}}}
	c.add(cacheItem{"one", 1})
	c.add(cacheItem{"two", 2})
	got, want := reflect.DeepEqual([]cacheItem{{"two", 2}}, c.items), true
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
