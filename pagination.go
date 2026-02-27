package hackeronecli

import (
	"net/url"
	"strconv"
)

type PageParams struct {
	Number int
	Size   int
}

func (p PageParams) Apply(params url.Values) url.Values {
	if params == nil {
		params = url.Values{}
	}
	if p.Number > 0 {
		params.Set("page[number]", strconv.Itoa(p.Number))
	}
	if p.Size > 0 {
		params.Set("page[size]", strconv.Itoa(p.Size))
	}
	return params
}
