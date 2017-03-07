package client

import "github.com/Sotera/go_watchman/loogo"
import "strings"

func main() {
	params := loogo.NewPagerParams{
		Params: loogo.QueryParams{},
		URL:    "http://localhost/api/events",
	}
	p, err := loogo.NewPager(params)
	if err != nil {
		panic(err)
	}

	c := client{}
	c.getNames(p)
}

type client struct{}

func (c client) getNames(p loogo.PagerInterface) (string, error) {

	docs, err := p.GetNext()
	if err != nil {
		return "", err
	}

	names := []string{}

	for _, d := range docs {
		names = append(names, d["name"].(string))
	}

	return strings.Join(names, " "), nil
}
