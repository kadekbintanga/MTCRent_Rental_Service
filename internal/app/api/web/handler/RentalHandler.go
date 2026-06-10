package handler

import (
	"net/http"
	motorcycleRepo "service/internal/motorcycle/repository"
	otherRepo "service/internal/other/repository"
	otherService "service/internal/other/service"
	form2 "service/internal/pkg/form"
	"service/internal/pkg/parser"
	"service/internal/rental/service"
	settingConfigRepo "service/internal/settingConfiguration/repository"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
	"github.com/globalxtreme/go-identifier/data"
	"github.com/gorilla/mux"
)

type RentalHandler struct{}

func (ctr RentalHandler) Create(w http.ResponseWriter, r *http.Request) {
	form := form2.RentalForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewRentalService()
	srv.SetEmployeeIdentifier(data.AuthEmployee(r))
	srv.SetCustomerRepository(otherRepo.NewCustomerRepository())
	srv.SetMotorcycleRepository(motorcycleRepo.NewMotorcycleRepository())
	srv.SetCustomerService(otherService.NewCustomerService())

	rental := srv.Create(form)

	psr := parser.RentalParser{Object: rental}
	res := xtremeres.Response{Object: psr.FirstFull()}
	res.Success(w)
}

func (ctr RentalHandler) Simulate(w http.ResponseWriter, r *http.Request) {
	form := form2.RentalSimulateForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewRentalService()

	simulate := srv.Simulate(mux.Vars(r)["uuid"], form)
	res := xtremeres.Response{Object: simulate}
	res.Success(w)
}

func (ctr RentalHandler) Update(w http.ResponseWriter, r *http.Request) {
	form := form2.RentalUpdateteForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewRentalService()
	rental := srv.Update(mux.Vars(r)["uuid"], form)

	psr := parser.RentalParser{Object: rental}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (ctr RentalHandler) Refund(w http.ResponseWriter, r *http.Request) {
	form := form2.RentalRefundForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewRentalService()
	refund := srv.Refund(mux.Vars(r)["uuid"], form)

	psr := parser.RentalRefundParser{Object: refund}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}

func (ctr RentalHandler) Return(w http.ResponseWriter, r *http.Request) {
	form := form2.RentalReturnForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewRentalService()
	srv.SetEmployeeIdentifier(data.AuthEmployee(r))
	srv.SetCustomerRepository(otherRepo.NewCustomerRepository())
	srv.SetSettingConfigurationRepository(settingConfigRepo.NewSettingConfigurationRepository())
	srv.SetMotorcycleRepository(motorcycleRepo.NewMotorcycleRepository())
	srv.SetCustomerService(otherService.NewCustomerService())

	rental := srv.Return(mux.Vars(r)["uuid"], form)
	psr := parser.RentalParser{Object: rental}
	res := xtremeres.Response{Object: psr.FirstFull()}
	res.Success(w)

}
