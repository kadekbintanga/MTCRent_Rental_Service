package form

import (
	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
)

type CustomerUpdateForm struct {
	ID        uint   `json:"id" validate:"required"`
	UUID      string `json:"uuid" validate:"required"`
	Name      string `json:"name" validate:"required"`
	IDNumber  string `json:"IDNumber" validate:"required"`
	SIMNumber string `json:"SIMNumber" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	StatusId  int    `json:"statusId" validate:"required"`
	Deleted   bool   `json:"deleted"`
}

func (rule *CustomerUpdateForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(rule)
}

func (rule *CustomerUpdateForm) AsyncWorkflowParse(payload interface{}) error {
	return core.BaseForm{}.AsyncWorkflowParse(payload, &rule)
}
