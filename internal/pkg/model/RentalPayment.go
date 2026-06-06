package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type RentalPayment struct {
	xtrememodel.BaseModel
	RentalId int     `gorm:"column:rentalId;type:int;not null"`
	Amount   float64 `gorm:"column:amount;type:float;not null"`
	MethodId int     `gorm:"column:methodId;type:int;not null"`
}

func (RentalPayment) TableName() string {
	return "rental_payments"
}

func (model RentalPayment) SetReference() uint {
	return model.BaseModel.ID
}
