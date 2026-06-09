package error

import (
	"net/http"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

func ErrXtremeRentalGet(internalMsg string) {
	xtremeres.Error(http.StatusNotFound, "Rental not found", internalMsg, false, nil)
}

func ErrXtremeRentalSave(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to save Rental", internalMsg, false, nil)
}

func ErrXtremeRentalDelete(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to delete Rental", internalMsg, false, nil)
}

func ErrXtremeRentalUpdate(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to update Rental", internalMsg, false, nil)
}

func ErrXtremeRentalSimulate(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to simulate Rental", internalMsg, false, nil)
}

func ErrXtremeRentalRefund(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to refund Rental", internalMsg, false, nil)
}

func ErrXtremeRentalReturn(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to return Rental", internalMsg, false, nil)
}

func ErrXtremeRentalInvalidPayload(internalMsg string) {
	xtremeres.Error(http.StatusBadRequest, "Invalid data request", internalMsg, false, nil)
}
