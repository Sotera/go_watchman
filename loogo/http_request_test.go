package loogo

import (
	"fmt"
	"testing"
)

func TestNewRequest(t *testing.T) {
	// HTTPRequester
	// doc:=Doc{}
	// err:=NewRequest(NewRequestParams{URL: url}, &doc)
	// url := "http://localhost:3000/api/jobsets/"
	// p1 := QueryParam{
	// 	QueryType: "Eq",
	// 	Field:     "featurizer",
	// 	Values:    []string{"image"},
	// }
	// p2 := QueryParam{
	// 	QueryType: "Eq",
	// 	Field:     "state",
	// 	Values:    []string{"done"},
	// }

	// params := QueryParams{
	// 	// p1,
	// }

	// obj := map[string]interface{}{
	// 	"start_time": 3,
	// 	"end_time":   4,
	// }

	// jobset, err := json.Marshal(obj)
	// if err != nil {
	// 	panic(err)
	// }
	// doc := &Doc{}
	// err = NewRequest(NewRequestParams{url, params, "POST", jobset}, doc)
	// if err != nil {
	// 	fmt.Println("err", err)
	// 	return
	// }
	// fmt.Println(doc)
}

func TestNewRequest2(t *testing.T) {
	url := "http://localhost:3000/api/jobsets/"
	p1 := QueryParam{
		QueryType: "Eq",
		Field:     "state",
		Values:    []string{"new"},
	}
	// p2 := QueryParam{
	// 	QueryType: "Eq",
	// 	Field:     "state",
	// 	Values:    []string{"done"},
	// }

	params := QueryParams{
		p1,
	}

	docs := &Docs{}
	err := NewRequest(NewRequestParams{url, params, "GET", nil}, docs)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	fmt.Println(docs)
}

func TestNewPager(t *testing.T) {
	// url := "http://localhost:3000/api/socialmediaposts/"
	// p1 := QueryParam{
	// 	QueryType: "Eq",
	// 	Field:     "featurizer",
	// 	Values:    []string{"image"},
	// }
	// // p2 := QueryParam{
	// // 	QueryType: "Eq",
	// // 	Field:     "state",
	// // 	Values:    []string{"done"},
	// // }

	// params := QueryParams{
	// 	p1,
	// }

	// pager, err := NewPager(NewPagerParams{
	// 	URL:      url,
	// 	Params:   params,
	// 	PageSize: 10,
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// docs, err := pager.GetNext()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// docs, _ = pager.GetNext()

	// var i = 1
	// docFunc := func(doc Doc, done func()) {
	// 	// fmt.Println(doc)
	// 	fmt.Println(i)
	// 	time.Sleep(0 * time.Second)
	// 	i++
	// 	done()
	// }
	// err = pager.PageOver(docFunc)
	// if err != nil {
	// 	fmt.Println(err)
	// 	// return
	// }

	// fmt.Println(len(docs), pager)
}
