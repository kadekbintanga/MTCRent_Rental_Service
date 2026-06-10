package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type SettingConfiguration struct {
	xtrememodel.BaseModel
	Key           string  `gorm:"column:key;type:varchar(250);not null"`
	Value         string  `gorm:"column:value;type:varchar(250);not null"`
	UpdatedBy     *string `gorm:"column:updatedBy;varchar(50);null"`
	UpdatedByName *string `gorm:"column:updatedByName;varchar(250);null"`
}

func (SettingConfiguration) TableName() string {
	return "setting_configurations"
}

func (model SettingConfiguration) SetReference() uint {
	return model.BaseModel.ID
}
