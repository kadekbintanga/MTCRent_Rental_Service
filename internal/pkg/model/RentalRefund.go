package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type RentalRefund struct {
	xtrememodel.BaseModel
	RentalId int     `gorm:"column:rentalId;type:int;not null"`
	Amount   float64 `gorm:"column:amount;type:float;not null"`
	MethodId int     `gorm:"column:methodId;type:int;not null"`
}

func (RentalRefund) TableName() string {
	return "rental_refunds"
}

func (model RentalRefund) SetReference() uint {
	return model.BaseModel.ID
}
