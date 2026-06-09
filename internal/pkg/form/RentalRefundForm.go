package form

import (
	"net/http"
	"service/internal/pkg/core"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
)

type RentalRefundForm struct {
	RefundMethodId int     `json:"refundMethodId"`
	RefundAmount   float64 `json:"refundAmount"`
}

func (rule *RentalRefundForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(rule)
}

func (rule *RentalRefundForm) APIParse(r *http.Request) {
	core.BaseForm{}.APIParse(r, &rule)
}
