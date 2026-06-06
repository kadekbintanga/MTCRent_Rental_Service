package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type SettingConfiguration struct {
	xtrememodel.BaseModel
	Key   string `gorm:"column:key;type:varchar(250);not null"`
	Value string `gorm:"column:value;type:varchar(250);not null"`
}

func (SettingConfiguration) TableName() string {
	return "setting_configurations"
}

func (model SettingConfiguration) SetReference() uint {
	return model.BaseModel.ID
}
