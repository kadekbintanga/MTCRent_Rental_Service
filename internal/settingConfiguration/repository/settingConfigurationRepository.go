package repository

import (
	"gorm.io/gorm"

	"service/internal/pkg/config"
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/form"
	"service/internal/pkg/model"
)

type SettingConfigurationRepository interface {
	core.TransactionInterface
	core.FirstRepository[form.SettingConfigurationFilterForm, model.SettingConfiguration]
	core.FindRepository[form.SettingConfigurationFilterForm, model.SettingConfiguration]

	Update(settingConfig model.SettingConfiguration, form form.SettingConfigurationForm) model.SettingConfiguration
}

func NewSettingConfigurationRepository(args ...*gorm.DB) SettingConfigurationRepository {
	repository := settingConfigurationRepository{}
	if len(args) > 0 {
		repository.Transaction = args[0]
	}

	return &repository
}

type settingConfigurationRepository struct {
	Transaction *gorm.DB
}

func (repo *settingConfigurationRepository) SetTransaction(tx *gorm.DB) {
	repo.Transaction = tx
}

func (repo *settingConfigurationRepository) FirstByForm(form form.SettingConfigurationFilterForm, args ...func(query *gorm.DB) *gorm.DB) model.SettingConfiguration {
	query := repo.prepareAndFilter(form)

	if len(args) > 0 {
		query = args[0](query)
	}

	var settingConfig model.SettingConfiguration
	err := query.First(&settingConfig).Error
	if err != nil {
		error2.ErrXtremeSettingConfigurationGet(err.Error())
	}
	return settingConfig
}

func (repo *settingConfigurationRepository) FindByForm(form form.SettingConfigurationFilterForm) []model.SettingConfiguration {
	query := repo.prepareAndFilter(form)

	var settingConfig []model.SettingConfiguration
	err := query.Find(&settingConfig).Error
	if err != nil {
		error2.ErrXtremeSettingConfigurationGet(err.Error())
	}
	return settingConfig
}

func (repo *settingConfigurationRepository) Update(settingConfig model.SettingConfiguration, form form.SettingConfigurationForm) model.SettingConfiguration {
	settingConfig.Value = form.Value

	err := repo.Transaction.Updates(&settingConfig).Error
	if err != nil {
		error2.ErrXtremeSettingConfigurationUpdate(err.Error())
	}
	return settingConfig
}

/** --- UNEXPORTED FUNCTIONS --- */
func (repo settingConfigurationRepository) prepareAndFilter(form form.SettingConfigurationFilterForm) *gorm.DB {
	query := config.PgSQL

	if form.ID > 0 {
		query = query.Where(`setting_configurations."id" = ?`, form.ID)
	}

	if form.Key != "" {
		query = query.Where(`setting_configurations."key" = ?`, form.Key)
	}

	return query
}
