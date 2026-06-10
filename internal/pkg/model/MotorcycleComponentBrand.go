package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type MotorcycleComponentBrand struct {
	xtrememodel.BaseModel
	Name          string  `gorm:"column:name;type:varchar(250);not null"`
	Default       bool    `gorm:"column:default"`
	CreatedBy     *string `gorm:"column:createdBy;varchar(50);null"`
	CreatedByName *string `gorm:"column:createdByName;varchar(250);null"`
	UpdatedBy     *string `gorm:"column:updatedBy;varchar(50);null"`
	UpdatedByName *string `gorm:"column:updatedByName;varchar(250);null"`

	Motorcycles []Motorcycle `gorm:"foreignKey:brandId"`
}

func (MotorcycleComponentBrand) TableName() string {
	return "motorcycle_component_brands"
}

func (model MotorcycleComponentBrand) SetReference() uint {
	return model.BaseModel.ID
}
