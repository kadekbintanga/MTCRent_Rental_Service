package form

import (
	"net/http"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"

	"service/internal/pkg/constant"
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
)

type MotorcycleStatusUpdateForm struct {
	StatusId int `json:"statusId" validate:"required"`
}

func (rule *MotorcycleStatusUpdateForm) Validate() {
	va := xtrememdw.Validator{}

	_, existStatusId := constant.MotorcycleStatus{}.OptionIDNames()[rule.StatusId]
	if !existStatusId {
		error2.ErrXtremeMotorcycleUpdate("Invalid status")
	}

	va.Make(rule)
}

func (rule *MotorcycleStatusUpdateForm) APIParse(r *http.Request) {
	core.BaseForm{}.APIParse(r, &rule)
}
