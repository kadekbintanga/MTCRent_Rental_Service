package handler

import (
	"net/http"
	form2 "service/internal/pkg/form"
	"service/internal/pkg/parser"
	"service/internal/rental/repository"
	"service/internal/rental/service"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
	"github.com/gorilla/mux"
)

type RentalRefundHandler struct{}

func (ctr RentalRefundHandler) Get(w http.ResponseWriter, r *http.Request) {
	form := form2.RentalRefundFilterForm{
		Preloads: []string{"Rental"},
	}

	form.FilterParse(r.URL.Query())

	repo := repository.NewRentalRefundRepository()
	refunds, pagination := repo.PaginateByForm(form)

	psr := parser.RentalRefundParser{Array: refunds}
	res := xtremeres.Response{Array: psr.Get(), Pagination: &pagination}
	res.Success(w)
}

func (ctr RentalRefundHandler) Create(w http.ResponseWriter, r *http.Request) {
	form := form2.RentalRefundForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewRentalRefundService()
	refund := srv.Create(mux.Vars(r)["uuid"], form)

	psr := parser.RentalRefundParser{Object: refund}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (ctr RentalRefundHandler) GetByRental(w http.ResponseWriter, r *http.Request) {
	form := form2.RentalRefundFilterForm{
		RentalUUID: mux.Vars(r)["uuid"],
	}

	form.FilterParse(r.URL.Query())

	repo := repository.NewRentalRefundRepository()
	refunds, pagination := repo.PaginateByForm(form)

	psr := parser.RentalRefundParser{Array: refunds}
	res := xtremeres.Response{Array: psr.Get(), Pagination: &pagination}
	res.Success(w)
}
