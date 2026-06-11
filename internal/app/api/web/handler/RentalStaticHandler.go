package handler

import (
	"net/http"
	"service/internal/pkg/constant"
	"service/internal/pkg/core"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

type RentalStaticHandler struct{}

func (ctr RentalStaticHandler) RentalStatus(w http.ResponseWriter, r *http.Request) {
	rentalStatus := core.IDName{}.Get(constant.RentalStatus{})

	res := xtremeres.Response{Array: rentalStatus}
	res.Success(w)
}

func (ctr RentalStaticHandler) RentalPaymentMethod(w http.ResponseWriter, r *http.Request) {
	paymentMethod := core.IDName{}.Get(constant.RentalPaymentMethod{})

	res := xtremeres.Response{Array: paymentMethod}
	res.Success(w)
}

func (ctr RentalStaticHandler) RentalRefundMethod(w http.ResponseWriter, r *http.Request) {
	refundMethod := core.IDName{}.Get(constant.RentalRefundMethod{})

	res := xtremeres.Response{Array: refundMethod}
	res.Success(w)
}
