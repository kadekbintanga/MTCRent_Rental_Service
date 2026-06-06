package form

import (
	"net/http"

	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"

	"service/internal/pkg/core"
)

type SettingConfigurationForm struct {
	Value string `json:"value" validate:"required"`
}

func (rule *SettingConfigurationForm) Validate() {
	va := xtrememdw.Validator{}
	va.Make(rule)
}

func (rule *SettingConfigurationForm) APIParse(r *http.Request) {
	core.BaseForm{}.APIParse(r, &rule)
}
