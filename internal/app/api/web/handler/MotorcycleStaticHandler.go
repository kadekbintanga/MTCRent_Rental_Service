package handler

import (
	"net/http"
	"service/internal/pkg/constant"
	"service/internal/pkg/core"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

type MotorcycleStaticHandler struct{}

func (ctr MotorcycleStaticHandler) MotorcycleType(w http.ResponseWriter, r *http.Request) {
	motorcycleType := core.IDName{}.Get(constant.MotorcycleType{})

	res := xtremeres.Response{Array: motorcycleType}
	res.Success(w)
}

func (ctr MotorcycleStaticHandler) MotorcycleStatus(w http.ResponseWriter, r *http.Request) {
	motorcycleStatus := core.IDName{}.Get(constant.MotorcycleStatus{})

	res := xtremeres.Response{Array: motorcycleStatus}
	res.Success(w)
}
