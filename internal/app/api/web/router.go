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

	var staticHandler handler.MotorcycleComponentStaticHandler
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
	router.HandleFunc("/{uuid}", motorcycleHandler.Delete).Methods("DELETE")
}

func settingRouter(router *mux.Router) {
	router = router.PathPrefix("/settings").Subrouter()

	var settingConfig handler.SettingConfigurationHandler
	router.HandleFunc("/configurations", settingConfig.Get).Methods("GET")
	router.HandleFunc("/configurations/{id}", settingConfig.Update).Methods("PUT")
}

func rentalRouter(router *mux.Router) {
	router = router.PathPrefix("/rentals").Subrouter()

	var staticHandler handler.RentalComponentStaticHandler
	router.HandleFunc("/components/statics/rental-statuses", staticHandler.RentalStatus).Methods("GET")
	router.HandleFunc("/components/statics/payment-methods", staticHandler.RentalPaymentMethod).Methods("GET")
	router.HandleFunc("/components/statics/refund-methods", staticHandler.RentalRefundMethod).Methods("GET")

	var paymentHandler handler.RentalPaymentHandler
	router.HandleFunc("/payments", paymentHandler.Get).Methods("GET")
	router.HandleFunc("/{uuid}/payments", paymentHandler.GetByRental).Methods("GET")
	router.HandleFunc("/{uuid}/payments", paymentHandler.Create).Methods("POST")

	var refundHandler handler.RentalRefundHandler
	router.HandleFunc("/refunds", refundHandler.Get).Methods("GET")
	router.HandleFunc("/{uuid}/refunds", refundHandler.GetByRental).Methods("GET")
	router.HandleFunc("/{uuid}/refunds", refundHandler.Create).Methods("POST")

	var rentalHandler handler.RentalHandler
	router.HandleFunc("", rentalHandler.Get).Methods("GET")
	router.HandleFunc("", rentalHandler.Create).Methods("POST")
	router.HandleFunc("/{uuid}", rentalHandler.Detail).Methods("GET")
	router.HandleFunc("/{uuid}", rentalHandler.Update).Methods("PUT")
	router.HandleFunc("/{uuid}/simulates", rentalHandler.Simulate).Methods("POST")
	router.HandleFunc("/{uuid}/returns", rentalHandler.Return).Methods("POST")
}
