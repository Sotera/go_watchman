package annotations

import (
	"github.com/Sotera/go_watchman/loogo"
)

type LoogoParserFactory interface {
	Generate() loogo.RequestParser
}

type ParserFactory struct {
}

func (pf ParserFactory) Generate() loogo.RequestParser {
	return &loogo.HTTPRequestParser{
		Client: &loogo.HTTPClient{},
	}
}
