package handler

import (
	"net/http"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
	"github.com/globalxtreme/go-identifier/data"
	"github.com/gorilla/mux"

	repository2 "service/internal/activity/repository"
	"service/internal/motorcycle/repository"
	"service/internal/motorcycle/service"
	"service/internal/pkg/core"
	form2 "service/internal/pkg/form"
	"service/internal/pkg/parser"
)

type MotorcycleComponentBrandHandler struct{}

func (ctr MotorcycleComponentBrandHandler) Get(w http.ResponseWriter, r *http.Request) {
	form := form2.MotorcycleComponentBrandFilterForm{
		Orders: map[string]string{"name": "ASC"},
	}

	form.FilterParse(r.URL.Query())

	repo := repository.NewMotorcycleComponentBrandRepository()
	motorcycleBrands, pagination := repo.PaginateByForm(form)

	psr := parser.MotorcycleBrandParser{Array: motorcycleBrands}

	res := xtremeres.Response{Array: psr.Get(), Pagination: &pagination}
	res.Success(w)
}

func (ctr MotorcycleComponentBrandHandler) Create(w http.ResponseWriter, r *http.Request) {
	form := form2.MotorcycleComponentBrandForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewMotorcycleComponentBrandService()
	srv.SetEmployeeIdentifier(data.AuthEmployee(r))
	srv.SetActivityRepository(repository2.NewActivityRepository())

	motorcycleBrand := srv.Create(form)

	psr := parser.MotorcycleBrandParser{Object: motorcycleBrand}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (ctr MotorcycleComponentBrandHandler) Update(w http.ResponseWriter, r *http.Request) {
	formFilter := form2.MotorcycleComponentBrandFilterForm{
		ID: core.ToInt(mux.Vars(r)["id"]),
	}
	form := form2.MotorcycleComponentBrandForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewMotorcycleComponentBrandService()
	srv.SetEmployeeIdentifier(data.AuthEmployee(r))
	srv.SetActivityRepository(repository2.NewActivityRepository())

	motorcycleBrand := srv.Update(formFilter, form)

	psr := parser.MotorcycleBrandParser{Object: motorcycleBrand}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)

}
