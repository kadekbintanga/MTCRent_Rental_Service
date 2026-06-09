package form

import (
	"net/url"

	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
)

type RentalFilterForm struct {
	ID             uint
	UUID           string
	CustomerUUID   string
	CustomerId     int
	MotorcycleUUID string
	StatusId       int
	Orders         map[string]string
	Preloads       []string
	Page           int
	Limit          int
}

func (f *RentalFilterForm) FilterParse(parameter url.Values) {
	if pageReq := parameter.Get("page"); pageReq != "" {
		f.Page = xtremepkg.ToInt(pageReq)
		f.Limit = xtremepkg.ToInt(parameter.Get("limit"))
	} else {
		f.Page = 1
	}
}
