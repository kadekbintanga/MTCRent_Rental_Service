package form

import (
	"net/http"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"

	"service/internal/pkg/constant"
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
)

type MotorcycleForm struct {
	PlateNumber string  `json:"plateNumber" validate:"required,max=10,alphanum"`
	Name        string  `json:"name" validate:"required,max=100"`
	TypeId      int     `json:"typeId" validate:"required"`
	Year        int     `json:"year" validate:"required"`
	PricePerDay float64 `json:"pricePerDay" validate:"required"`
	StatusId    int     `json:"statusId" validate:"required"`
	BrandId     int     `json:"brandId" validate:"required"`
}

func (rule *MotorcycleForm) Validate() {
	va := xtrememdw.Validator{}
	_, existTypeId := constant.MotorcycleType{}.OptionIDNames()[rule.TypeId]
	if !existTypeId {
		error2.ErrXtremeMotorcycleSave("Invalid type")
	}

	_, existStatusId := constant.MotorcycleStatus{}.OptionIDNames()[rule.StatusId]
	if !existStatusId {
		error2.ErrXtremeMotorcycleSave("Invalid status")
	}

	va.Make(rule)
}

func (rule *MotorcycleForm) APIParse(r *http.Request) {
	core.BaseForm{}.APIParse(r, &rule)
}
