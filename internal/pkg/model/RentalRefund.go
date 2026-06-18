package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type RentalRefund struct {
	xtrememodel.BaseModel
	Number        string  `gorm:"column:number;type:varchar(250);not null"`
	RentalId      uint    `gorm:"column:rentalId;type:int;not null"`
	Amount        float64 `gorm:"column:amount;type:float;not null"`
	MethodId      int     `gorm:"column:methodId;type:int;not null"`
	CreatedBy     *string `gorm:"column:createdBy;varchar(50);null"`
	CreatedByName *string `gorm:"column:createdByName;varchar(250);null"`

	Rental Rental `gorm:"foreignKey:RentalId"`
}

func (RentalRefund) TableName() string {
	return "rental_refunds"
}

func (model RentalRefund) SetReference() uint {
	return model.BaseModel.ID
}
