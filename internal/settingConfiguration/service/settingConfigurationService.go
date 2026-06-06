package service

import (
	"fmt"

	"github.com/globalxtreme/go-identifier/data"
	"gorm.io/gorm"

	"service/internal/pkg/activity"
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	"service/internal/pkg/form"
	form2 "service/internal/pkg/form"
	"service/internal/pkg/model"
	"service/internal/pkg/parser"
	"service/internal/settingConfiguration/repository"
)

type SettingConfigurationService interface {
	SetTransaction(tx *gorm.DB)
	SetEmployeeIdentifier(employee data.EmployeeIdentifierData)

	Update(filterForm form2.SettingConfigurationFilterForm, form form.SettingConfigurationForm) model.SettingConfiguration
}

func NewSettingConfigurationService() SettingConfigurationService {
	return &settingConfigurationService{}
}

type settingConfigurationService struct {
	tx         *gorm.DB
	repository repository.SettingConfigurationRepository
	employee   data.EmployeeIdentifierData
}

func (srv *settingConfigurationService) SetTransaction(tx *gorm.DB) {
	srv.tx = tx
}

func (srv *settingConfigurationService) SetEmployeeIdentifier(employee data.EmployeeIdentifierData) {
	srv.employee = employee
}

func (srv *settingConfigurationService) Update(filterForm form2.SettingConfigurationFilterForm, form form.SettingConfigurationForm) model.SettingConfiguration {
	srv.repository = repository.NewSettingConfigurationRepository()
	settingConfig := srv.repository.FirstByForm(filterForm)

	parser := parser.SettingConfigurationParser{Object: settingConfig}

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)
		useActivity := activity.UseActivity{Employee: srv.employee}.SetReference(&settingConfig).SetParser(&parser).SetOldProperty(constant.ACTION_UPDATE)

		settingConfig = srv.repository.Update(settingConfig, form)

		parser.Object = settingConfig
		useActivity.SetReference(&settingConfig).SetParser(&parser).SetNewProperty(constant.ACTION_UPDATE).
			Save(fmt.Sprintf("Update Setting Configuration %s [%d]", settingConfig.Key, settingConfig.ID))
		return nil
	})
	return settingConfig

}
