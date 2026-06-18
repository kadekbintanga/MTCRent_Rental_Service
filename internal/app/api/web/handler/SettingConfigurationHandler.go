package handler

import (
	"net/http"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
	"github.com/globalxtreme/go-identifier/data"
	"github.com/gorilla/mux"

	"service/internal/pkg/core"
	form2 "service/internal/pkg/form"
	"service/internal/pkg/parser"
	"service/internal/setting/repository"
	"service/internal/setting/service"
)

type SettingConfigurationHandler struct{}

func (ctr SettingConfigurationHandler) Get(w http.ResponseWriter, r *http.Request) {
	form := form2.SettingConfigurationFilterForm{}

	repo := repository.NewSettingConfigurationRepository()
	settingConfig := repo.FindByForm(form)

	psr := parser.SettingConfigurationParser{Array: settingConfig}

	res := xtremeres.Response{Array: psr.Get()}
	res.Success(w)
}

func (ctr SettingConfigurationHandler) Update(w http.ResponseWriter, r *http.Request) {
	formFilter := form2.SettingConfigurationFilterForm{
		ID: core.ToInt(mux.Vars(r)["id"]),
	}
	form := form2.SettingConfigurationForm{}
	form.APIParse(r)
	form.Validate()

	srv := service.NewSettingConfigurationService()
	srv.SetEmployeeIdentifier(data.AuthEmployee(r))

	settingConfig := srv.Update(formFilter, form)

	psr := parser.SettingConfigurationParser{Object: settingConfig}
	res := xtremeres.Response{Object: psr.First()}
	res.Success(w)
}
