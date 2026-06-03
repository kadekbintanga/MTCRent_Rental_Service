package form

import (
	"net/url"

	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
)

type MotorcycleComponentBrandFilterForm struct {
	ID     int
	Name   string
	Search string
	Orders map[string]string
	Page   int
	Limit  int
}

func (f *MotorcycleComponentBrandFilterForm) FilterParse(parameter url.Values) {
	if searchReq := parameter.Get("search"); len(searchReq) >= 3 {
		f.Search = searchReq
	}

	if pageReq := parameter.Get("page"); pageReq != "" {
		f.Page = xtremepkg.ToInt(pageReq)
		f.Limit = xtremepkg.ToInt(parameter.Get("limit"))
	} else {
		f.Page = 1
	}
}
