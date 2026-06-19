package handler

import (
	"net/http"
	form2 "service/internal/pkg/form"
	"service/internal/pkg/parser"
	"service/internal/rental/repository"
	"service/internal/rental/service"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
	"github.com/globalxtreme/go-identifier/data"
	"github.com/gorilla/mux"
)

type RentalPaymentHandler struct{}

func (ctr RentalPaymentHandler) Get(w http.ResponseWriter, r *http.Request) {
	form := form2.RentalPaymentFilterForm{
		Preloads: []string{"Rental"},
	}

	form.FilterParse(r.URL.Query())

	repo := repository.NewRentalPaymentRepository()
	payments, pagination := repo.PaginateByForm(form)

	psr := parser.RentalPaymentParser{Array: payments}
	res := xtremeres.Response{Array: psr.Get(), Pagination: &pagination}
	res.Success(w)
}

func (ctr RentalPaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	form := form2.RentalPaymentForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewRentalPaymentService()
	srv.SetEmployeeIdentifier(data.AuthEmployee(r))
	payment := srv.Create(mux.Vars(r)["uuid"], form)

	psr := parser.RentalPaymentParser{Object: payment}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (ctr RentalPaymentHandler) GetByRental(w http.ResponseWriter, r *http.Request) {
	form := form2.RentalPaymentFilterForm{
		RentalUUID: mux.Vars(r)["uuid"],
	}

	form.FilterParse(r.URL.Query())

	repo := repository.NewRentalPaymentRepository()
	payments, pagination := repo.PaginateByForm(form)

	psr := parser.RentalPaymentParser{Array: payments}
	res := xtremeres.Response{Array: psr.Get(), Pagination: &pagination}
	res.Success(w)
}
