package web

import (
	"github.com/gorilla/mux"

	"service/internal/app/api/web/handler"
)

func Register(router *mux.Router) {
	activityRouter(router)
	testingRouter(router) // TODO: Hanya contoh. nanti langsung hapus saja
	motorcycleRouter(router)
	settingRouter(router)
	rentalPaymentRouter(router)
	rentalRefundRouter(router)
	rentalRouter(router)
}

func activityRouter(router *mux.Router) {
	router.HandleFunc("/activities", handler.ActivityHandler{}.Get).Methods("GET")
}

// TODO: Hanya contoh. nanti langsung hapus saja
func testingRouter(router *mux.Router) {
	var testingHandler handler.TestingHandler
	router.HandleFunc("/testings", testingHandler.Get).Methods("GET")
	router.HandleFunc("/testings", testingHandler.Create).Methods("POST")
	router.HandleFunc("/testings/upload/file", testingHandler.UploadByFile).Methods("POST")
	router.HandleFunc("/testings/upload/content", testingHandler.UploadByContent).Methods("POST")
}

func motorcycleRouter(router *mux.Router) {
	router = router.PathPrefix("/motorcycles").Subrouter()

	var staticHandler handler.MotorcycleStaticHandler
	router.HandleFunc("/components/statics/motorcycle-types", staticHandler.MotorcycleType).Methods("GET")
	router.HandleFunc("/components/statics/motorcycle-statuses", staticHandler.MotorcycleStatus).Methods("GET")

	var brandHandler handler.MotorcycleComponentBrandHandler
	router.HandleFunc("/components/brands", brandHandler.Get).Methods("GET")
	router.HandleFunc("/components/brands", brandHandler.Create).Methods("POST")
	router.HandleFunc("/components/brands/{id}", brandHandler.Update).Methods("PUT")
	router.HandleFunc("/components/brands/{id}", brandHandler.Delete).Methods("DELETE")

	var motorcycleHandler handler.MotorcycleHandler
	router.HandleFunc("", motorcycleHandler.Get).Methods("GET")
	router.HandleFunc("", motorcycleHandler.Create).Methods("POST")
	router.HandleFunc("/{uuid}", motorcycleHandler.Detail).Methods("GET")
	router.HandleFunc("/{uuid}", motorcycleHandler.Update).Methods("PUT")
}

func settingRouter(router *mux.Router) {
	router = router.PathPrefix("/settings").Subrouter()

	var settingConfig handler.SettingConfigurationHandler
	router.HandleFunc("/configurations", settingConfig.Get).Methods("GET")
	router.HandleFunc("/configurations/{id}", settingConfig.Update).Methods("PUT")
}

func rentalPaymentRouter(router *mux.Router) {
	router = router.PathPrefix("/payments").Subrouter()

	var paymentHandler handler.RentalPaymentHandler
	router.HandleFunc("", paymentHandler.Get).Methods("GET")
}

func rentalRefundRouter(router *mux.Router) {
	router = router.PathPrefix("/refunds").Subrouter()

	var refundHandler handler.RentalRefundHandler
	router.HandleFunc("", refundHandler.Get).Methods("GET")
}

func rentalRouter(router *mux.Router) {
	var staticHandler handler.RentalStaticHandler
	router.HandleFunc("/components/statics/rental-statuses", staticHandler.RentalStatus).Methods("GET")
	router.HandleFunc("/components/statics/payment-methods", staticHandler.RentalPaymentMethod).Methods("GET")
	router.HandleFunc("/components/statics/refund-methods", staticHandler.RentalRefundMethod).Methods("GET")

	var rentalHandler handler.RentalHandler
	router.HandleFunc("", rentalHandler.Get).Methods("GET")
	router.HandleFunc("", rentalHandler.Create).Methods("POST")
	router.HandleFunc("/{uuid}", rentalHandler.Detail).Methods("GET")
	router.HandleFunc("/{uuid}", rentalHandler.Update).Methods("PUT")
	router.HandleFunc("/{uuid}/simulates", rentalHandler.Simulate).Methods("POST")
	router.HandleFunc("/{uuid}/refunds", rentalHandler.Refund).Methods("POST")
	router.HandleFunc("/{uuid}/return", rentalHandler.Return).Methods("POST")
}
