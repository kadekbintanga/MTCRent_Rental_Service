package error

import (
	"net/http"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

func ErrXtremeSettingConfigurationGet(internalMsg string) {
	xtremeres.Error(http.StatusNotFound, "Motorcycle Brand not found", internalMsg, false, nil)
}

func ErrXtremeSettingConfigurationUpdate(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to update Motorcycle Brand", internalMsg, false, nil)
}
