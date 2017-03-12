package client

import "github.com/Sotera/go_watchman/loogo"
import "strings"

func main() {
	// code in here simply creates a 'real' pager and calls 'getNames'.
	// nothing interesting to test here.
	params := loogo.NewPagerParams{
		Params: loogo.QueryParams{},
		URL:    "http://localhost:3000/api/events",
	}
	p, err := loogo.NewPager(params)
	if err != nil {
		panic(err)
	}

	c := client{
		pager: p,
	}
	c.getNames()
}

type client struct {
	pager loogo.PagerInterface
}

// this is a function we want to test and can do so by providing our own
// custom pager.
func (c client) getNames() (string, error) {

	docs, err := c.pager.GetNext()
	if err != nil {
		return "", err
	}

	names := []string{}

	for _, d := range docs {
		names = append(names, d["name"].(string))
	}

	return strings.Join(names, " "), nil
}
