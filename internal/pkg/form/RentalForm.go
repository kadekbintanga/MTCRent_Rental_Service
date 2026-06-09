package form

import (
	"net/http"
	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
)

type RentalForm struct {
	CustomerUUID    string  `json:"customerUUID" validate:"required"`
	MotorcycleUUID  string  `json:"motorcycleUUID" validate:"required"`
	ReturnDatePlan  string  `json:"returnDatePlan" validate:"required"`
	PaymentMethodId int     `json:"paymentMethodId"`
	PaymentAmount   float64 `json:"paymentAmount"`
}

func (rule *RentalForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(rule)
}

func (rule *RentalForm) APIParse(r *http.Request) {
	core.BaseForm{}.APIParse(r, &rule)
}
