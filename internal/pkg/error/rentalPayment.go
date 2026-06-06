package error

import (
	"net/http"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

func ErrXtremeRentalPaymentGet(internalMsg string) {
	xtremeres.Error(http.StatusNotFound, "Rental Payment not found", internalMsg, false, nil)
}

func ErrXtremeRentalPaymentSave(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to save Rental Payment", internalMsg, false, nil)
}

func ErrXtremeRentalPaymentDelete(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to delete Rental Payment", internalMsg, false, nil)
}

func ErrXtremeRentalPaymentUpdate(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to update Rental Payment", internalMsg, false, nil)
}
