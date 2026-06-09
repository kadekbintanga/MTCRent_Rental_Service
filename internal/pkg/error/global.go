package error

import (
	"net/http"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

func ErrXtremePrivateAPIAuthentication(internalMsg string) {
	xtremeres.Error(http.StatusUnauthorized, "Doesn't have access to private api", internalMsg, false, nil)
}

func ErrXtremeAsyncWorkflowPush(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to push async workflow push", internalMsg, false, nil)
}

func ErrXtremePermission() {
	xtremeres.Error(http.StatusForbidden, "Permission restricted", "", false, nil)
}

func ErrXtremeRole() {
	xtremeres.Error(http.StatusForbidden, "Role restricted", "", false, nil)
}

func ErrXtremeINCRNumber(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to generate number", internalMsg, false, nil)
}

func ErrXtremeInvalidPayload(internalMsg string) {
	xtremeres.Error(http.StatusBadRequest, "Invalid data request", internalMsg, false, nil)
}
