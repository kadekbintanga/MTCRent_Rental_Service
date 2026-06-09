package form

import (
	"net/http"
	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
)

type RentalSimulateForm struct {
	ReturnDate string `json:"returnDate" validate:"required"`
}

func (rule *RentalSimulateForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(rule)
}

func (rule *RentalSimulateForm) APIParse(r *http.Request) {
	core.BaseForm{}.APIParse(r, &rule)
}
