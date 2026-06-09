package form

import (
	"net/http"
	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
)

type RentalReturnForm struct {
	ReturnDateActual string  `json:"returnDateActual" validate:"required"`
	PaymentMethodId  int     `json:"paymentMethodId"`
	PaymentAmount    float64 `json:"paymentAmount"`
	Note             string  `json:"note"`
}

func (rule *RentalReturnForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(rule)
}

func (rule *RentalReturnForm) APIParse(r *http.Request) {
	core.BaseForm{}.APIParse(r, &rule)
}
