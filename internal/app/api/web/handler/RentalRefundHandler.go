package handler

import (
	"net/http"
	form2 "service/internal/pkg/form"
	"service/internal/pkg/parser"
	"service/internal/rental/repository"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
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
