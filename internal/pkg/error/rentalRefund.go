package error

import (
	"net/http"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

func ErrXtremeRentalRefundGet(internalMsg string) {
	xtremeres.Error(http.StatusNotFound, "Rental Refund not found", internalMsg, false, nil)
}

func ErrXtremeRentalRefundSave(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to save Rental Refund", internalMsg, false, nil)
}

func ErrXtremeRentalRefundDelete(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to delete Rental Refund", internalMsg, false, nil)
}

func ErrXtremeRentalRefundUpdate(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to update Rental Refund", internalMsg, false, nil)
}
