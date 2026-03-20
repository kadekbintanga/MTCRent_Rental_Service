package form

import (
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"net/url"
)

type TestingFilterForm struct {
	ID             uint
	FromDate       string
	ToDate         string
	Search         string
	Preloads       []string
	Orders         map[string]string
	UseTransaction bool
	WithPagination bool
	Page           int
	Limit          int
}

func (f *TestingFilterForm) FilterParse(parameter url.Values) {
	f.FromDate = parameter.Get("fromDate")
	f.ToDate = parameter.Get("toDate")

	if searchReq := parameter.Get("search"); len(searchReq) >= 3 {
		f.Search = searchReq
	}

	f.WithPagination = true
	if pageReq := parameter.Get("page"); pageReq != "" {
		f.Page = xtremepkg.ToInt(pageReq)
		f.Limit = xtremepkg.ToInt(parameter.Get("limit"))
	}
}
