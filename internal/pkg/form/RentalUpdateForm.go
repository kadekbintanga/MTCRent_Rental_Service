package form

import (
	"net/http"
	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
)

type RentalUpdateteForm struct {
	ReturnDatePlan string `json:"returnDatePlan" validate:"required"`
	Note           string `json:"note"`
}

func (rule *RentalUpdateteForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(rule)
}

func (rule *RentalUpdateteForm) APIParse(r *http.Request) {
	core.BaseForm{}.APIParse(r, &rule)
}
