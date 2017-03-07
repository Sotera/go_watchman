package go_watchman

import (
	"github.com/sotera/go_watchman/loogo"
)

type pager_factory struct{}

func (pf pager_factory) generate(params loogo.NewPagerParams) (loogo.PagerInterface, error) {
	return loogo.NewPager(params)
}
