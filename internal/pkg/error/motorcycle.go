package error

import (
	"net/http"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

// TODO: Hanya contoh. nanti langsung hapus saja

func ErrXtremeMotorcycleGet(internalMsg string) {
	xtremeres.Error(http.StatusNotFound, "Motorcycle  not found", internalMsg, false, nil)
}

func ErrXtremeMotorcycleSave(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to save Motorcycle ", internalMsg, false, nil)
}

func ErrXtremeMotorcycleDelete(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to delete Motorcycle ", internalMsg, false, nil)
}

func ErrXtremeMotorcycleUpdate(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to update Motorcycle ", internalMsg, false, nil)
}
