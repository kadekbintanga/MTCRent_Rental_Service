package form

import (
	"net/url"

	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
)

type RentalRefundFilterForm struct {
	ID         uint
	RentalUUID string
	MethodId   int
	Preloads   []string
	Page       int
	Limit      int
}

func (f *RentalRefundFilterForm) FilterParse(parameter url.Values) {
	if pageReq := parameter.Get("page"); pageReq != "" {
		f.Page = xtremepkg.ToInt(pageReq)
		f.Limit = xtremepkg.ToInt(parameter.Get("limit"))
	} else {
		f.Page = 1
	}
}
