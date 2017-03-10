package loogo

import (
	"fmt"
	_ "reflect"
	"testing"
	"time"
)

func TestNewPager(t *testing.T) {
	url := "http://localhost:3000/api/socialmediaposts/"
	p1 := QueryParam{
		QueryType: "Eq",
		Field:     "featurizer",
		Values:    []string{"image"},
	}
	// p2 := QueryParam{
	// 	QueryType: "Eq",
	// 	Field:     "state",
	// 	Values:    []string{"done"},
	// }

	params := QueryParams{
		p1,
	}

	pager, err := NewPager(NewPagerParams{
		URL:      url,
		Params:   params,
		PageSize: 10,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	// docs, err := pager.GetNext()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// docs, _ = pager.GetNext()

	var i = 1
	docFunc := func(doc Doc, done func()) {
		// fmt.Println(doc)
		fmt.Println(i)
		time.Sleep(0 * time.Second)
		i++
		done()
	}
	err = pager.PageOver(docFunc)
	if err != nil {
		fmt.Println(err)
		// return
	}

	// fmt.Println(len(docs), pager)
	fmt.Println(pager)
}
