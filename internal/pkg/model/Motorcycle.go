package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type Motorcycle struct {
	xtrememodel.BaseModelUUID
	PlateNumber string  `gorm:"column:plateNumber;type:varchar(250);not null"`
	Name        string  `gorm:"column:name;type:varchar(250);not null"`
	TypeId      int     `gorm:"column:typeId;type:int"`
	Year        int     `gorm:"column:year;type:int"`
	PricePerDay float64 `gorm:"column:pricePerDay;type:float"`
	StatusId    int     `gorm:"column:statusId;type:int"`
	BrandId     int     `gorm:"column:brandId;type:int"`

	Brand MotorcycleComponentBrand `gorm:"foreignKey:brandId"`
}

func (Motorcycle) TableName() string {
	return "motorcycles"
}

func (model Motorcycle) SetReference() uint {
	return model.BaseModelUUID.ID
}
