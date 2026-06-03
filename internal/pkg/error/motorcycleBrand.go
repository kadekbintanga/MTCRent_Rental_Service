package error

import (
	"net/http"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

// TODO: Hanya contoh. nanti langsung hapus saja

func ErrXtremeMotorcycleBrandGet(internalMsg string) {
	xtremeres.Error(http.StatusNotFound, "Motorcycle Brand not found", internalMsg, false, nil)
}

func ErrXtremeMotorcycleBrandSave(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to save Motorcycle Brand", internalMsg, false, nil)
}

func ErrXtremeMotorcycleBrandDelete(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to delete Motorcycle Brand", internalMsg, false, nil)
}

func ErrXtremeMotorcycleBrandUpdate(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to update Motorcycle Brand", internalMsg, false, nil)
}
