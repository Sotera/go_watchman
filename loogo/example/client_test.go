package client

import (
	"testing"

	"github.com/Sotera/go_watchman/loogo"
)

type CustomPager struct{}

func (p *CustomPager) GetNext() (loogo.Docs, error) {
	return loogo.Docs{
		loogo.Doc{
			"id":   "123",
			"name": "doggy",
		},
		loogo.Doc{
			"id":   "456",
			"name": "bear",
		},
	}, nil
}

func (p *CustomPager) PageOver(docFunc func(doc loogo.Doc, done func())) error {
	return nil
}

// use simple TestPager from loogo
func Test_client_doWork_testPager(t *testing.T) {
	pager := &loogo.TestPager{}

	c := client{}

	want := "apple pear"
	if got, _ := c.getNames(pager); got != want {
		t.Errorf("getNames() = %v, want %v", got, want)
	}
}

// use custom pager
func Test_client_doWork_customPager(t *testing.T) {
	pager := &CustomPager{}

	c := client{}

	want := "doggy bear"
	if got, _ := c.getNames(pager); got != want {
		t.Errorf("getNames() = %v, want %v", got, want)
	}
}
