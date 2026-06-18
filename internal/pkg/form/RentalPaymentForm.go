package form

import (
	"net/http"
	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
)

type RentalPaymentForm struct {
	PaymentMethodId int     `json:"paymentMethodId"`
	PaymentAmount   float64 `json:"paymentAmount"`
}

func (rule *RentalPaymentForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(rule)
}

func (rule *RentalPaymentForm) APIParse(r *http.Request) {
	core.BaseForm{}.APIParse(r, &rule)
}
