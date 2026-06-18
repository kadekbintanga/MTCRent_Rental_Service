package handler

import (
	"net/http"
	"service/internal/pkg/constant"
	"service/internal/pkg/core"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

type RentalComponentStaticHandler struct{}

func (ctr RentalComponentStaticHandler) RentalStatus(w http.ResponseWriter, r *http.Request) {
	rentalStatus := core.IDName{}.Get(constant.RentalStatus{})

	res := xtremeres.Response{Array: rentalStatus}
	res.Success(w)
}

func (ctr RentalComponentStaticHandler) RentalPaymentMethod(w http.ResponseWriter, r *http.Request) {
	paymentMethod := core.IDName{}.Get(constant.RentalPaymentMethod{})

	res := xtremeres.Response{Array: paymentMethod}
	res.Success(w)
}

func (ctr RentalComponentStaticHandler) RentalRefundMethod(w http.ResponseWriter, r *http.Request) {
	refundMethod := core.IDName{}.Get(constant.RentalRefundMethod{})

	res := xtremeres.Response{Array: refundMethod}
	res.Success(w)
}
