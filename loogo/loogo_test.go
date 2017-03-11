package loogo

import "testing"

type TestRequestParser struct {
}

func (p *TestRequestParser) NewRequest(params NewRequestParams, result interface{}) error {
	return nil
}

func TestGetNextHasData(t *testing.T) {
	pager := &Pager{Parser: &TestRequestParser{}, TotalCount: 42}
	pager.GetNext()
	pager.GetNext()
	got := pager.CurrentPage
	want := 2
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetNextNoData(t *testing.T) {
	pager := &Pager{Parser: &TestRequestParser{}, TotalCount: 0}
	pager.GetNext()
	pager.GetNext()
	got := pager.CurrentPage
	want := 0
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestPageOver(t *testing.T) {
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
	// // docs, err := pager.GetNext()
	// // if err != nil {
	// // 	fmt.Println(err)
	// // 	return
	// // }
	// // docs, _ = pager.GetNext()

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

	// // fmt.Println(len(docs), pager)
	// fmt.Println(pager)
}
