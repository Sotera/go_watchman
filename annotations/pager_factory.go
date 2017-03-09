package annotations

import (
	"github.com/Sotera/go_watchman/loogo"
)

type LoogoPagerFactory interface {
	Generate(params loogo.NewPagerParams) (loogo.PagerInterface, error)
}

type PagerFactory struct {
}

func (pf PagerFactory) Generate(params loogo.NewPagerParams) (loogo.PagerInterface, error) {
	return loogo.NewPager(params)
}
