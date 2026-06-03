package form

import (
	"net/http"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"

	"service/internal/pkg/core"
)

type MotorcycleComponentBrandForm struct {
	Name string `json:"name" validate:"required"`
}

func (rule *MotorcycleComponentBrandForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(rule)
}

func (rule *MotorcycleComponentBrandForm) APIParse(r *http.Request) {
	core.BaseForm{}.APIParse(r, &rule)
}
