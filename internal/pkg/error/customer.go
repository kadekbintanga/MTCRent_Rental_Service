package error

import (
	"net/http"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

func ErrXtremeCustomerGet(internalMsg string) {
	xtremeres.Error(http.StatusNotFound, "Customer not found", internalMsg, false, nil)
}

func ErrXtremeCustomerSave(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to save Customer", internalMsg, false, nil)
}

func ErrXtremeCustomerDelete(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to delete Customer", internalMsg, false, nil)
}

func ErrXtremeCustomerUpdate(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to update Customer", internalMsg, false, nil)
}
