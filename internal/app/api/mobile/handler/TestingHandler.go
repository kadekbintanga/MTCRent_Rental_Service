package handler

import (
	xtremeres "github.com/globalxtreme/go-core/v2/response"
	"net/http"
	form2 "service/internal/pkg/form"
	"service/internal/pkg/parser"
	"service/internal/testing/repository"
)

type TestingHandler struct{}

func (ctr TestingHandler) Get(w http.ResponseWriter, r *http.Request) {
	form := form2.TestingFilterForm{
		Preloads: []string{"Subs"},
		Orders:   map[string]string{"id": "DESC"},
	}
	form.FilterParse(r.URL.Query())

	repo := repository.NewTestingRepository()
	testings, pagination := repo.PaginateByForm(form)

	psr := parser.TestingParser{Array: testings}

	res := xtremeres.Response{Array: psr.Get(), Pagination: &pagination}
	res.Success(w)
}
