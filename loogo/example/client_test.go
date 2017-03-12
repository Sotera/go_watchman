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
func TestClientDoWorkTestPager(t *testing.T) {
	pager := &loogo.TestPager{}

	c := client{
		pager: pager,
	}

	want := "apple pear"
	if got, _ := c.getNames(); got != want {
		t.Errorf("getNames() = %v, want %v", got, want)
	}
}

// use custom pager
func TestClientDoWorkCustomPager(t *testing.T) {
	pager := &CustomPager{}

	c := client{
		pager: pager,
	}

	want := "doggy bear"
	if got, _ := c.getNames(); got != want {
		t.Errorf("getNames() = %v, want %v", got, want)
	}
}
