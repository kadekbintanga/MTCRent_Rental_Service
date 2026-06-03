package handler

import (
	"net/http"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
	"github.com/globalxtreme/go-identifier/data"
	"github.com/gorilla/mux"

	repository2 "service/internal/activity/repository"
	"service/internal/motorcycle/repository"
	"service/internal/motorcycle/service"
	form2 "service/internal/pkg/form"
	"service/internal/pkg/parser"
)

type MotorcycleHandler struct{}

func (ctr MotorcycleHandler) Get(w http.ResponseWriter, r *http.Request) {
	form := form2.MotorcycleFilterForm{
		Preloads: []string{"Brand"},
	}

	form.FilterParse(r.URL.Query())

	repo := repository.NewMotorcycleRepository()
	motorcycles, pagination := repo.PaginateByForm(form)

	psr := parser.MotorcycleParser{Array: motorcycles}

	res := xtremeres.Response{Array: psr.Get(), Pagination: &pagination}
	res.Success(w)
}

func (ctr MotorcycleHandler) Detail(w http.ResponseWriter, r *http.Request) {
	form := form2.MotorcycleFilterForm{
		UUID:     mux.Vars(r)["uuid"],
		Preloads: []string{"Brand"},
	}

	repo := repository.NewMotorcycleRepository()
	motorcycle := repo.FirstByForm(form)

	psr := parser.MotorcycleParser{Object: motorcycle}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (ctr MotorcycleHandler) Create(w http.ResponseWriter, r *http.Request) {
	form := form2.MotorcycleForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewMotorcycleService()
	srv.SetEmployeeIdentifier(data.AuthEmployee(r))
	srv.SetActivityRepository(repository2.NewActivityRepository())

	motorcycle := srv.Create(form)

	psr := parser.MotorcycleParser{Object: motorcycle}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (ctr MotorcycleHandler) Update(w http.ResponseWriter, r *http.Request) {
	formFilter := form2.MotorcycleFilterForm{
		UUID: mux.Vars(r)["uuid"],
	}
	form := form2.MotorcycleForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewMotorcycleService()
	srv.SetEmployeeIdentifier(data.AuthEmployee(r))
	srv.SetActivityRepository(repository2.NewActivityRepository())

	motorcycle := srv.Update(formFilter, form)

	psr := parser.MotorcycleParser{Object: motorcycle}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}
