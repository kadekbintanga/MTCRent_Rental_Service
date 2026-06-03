package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type MotorcycleComponentBrand struct {
	xtrememodel.BaseModel
	Name string `gorm:"column:name;type:varchar(250);not null"`

	Motorcycles []Motorcycle `gorm:"foreignKey:brandId"`
}

func (MotorcycleComponentBrand) TableName() string {
	return "motorcycle_component_brands"
}

func (model MotorcycleComponentBrand) SetReference() uint {
	return model.BaseModel.ID
}
