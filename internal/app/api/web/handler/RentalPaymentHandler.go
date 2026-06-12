package handler

import (
	"net/http"
	form2 "service/internal/pkg/form"
	"service/internal/pkg/parser"
	"service/internal/rental/repository"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
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
