package model

import (
	"time"

	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type Rental struct {
	xtrememodel.BaseModelUUID
	CustomerId            int       `gorm:"column:customerId;type:int;not null"`
	MotorcycleId          int       `gorm:"column:motorcycleId;type:int;not null"`
	MotorcyclePlateNumber string    `gorm:"column:motorcyclePlateNumber;type:varchar(250);not null"`
	RentDate              time.Time `gorm:"column:rentDate;type:date;not null"`
	ReturnDatePlan        time.Time `gorm:"column:returnDatePlan;type:date;not null"`
	ReturnDateActual      time.Time `gorm:"column:returnDateActual;type:date"`
	PricePerDay           float64   `gorm:"column:pricePerDay;type:float;not null"`
	RentDay               int       `gorm:"column:rentDay;type:int"`
	TotalRentPrice        float64   `gorm:"column:totalRentPrice;type:float"`
	LateDay               int       `gorm:"column:lateDat;type:int"`
	PinaltyPrice          float64   `gorm:"column:pinaltyPrice;type:float"`
	StatusId              int       `gorm:"column:statusId;type:int;not null"`
	Note                  string    `gorm:"column:note;type:text"`

	Customer      Customer      `gorm:"foreignKey:customerId"`
	Motorcycle    Motorcycle    `gorm:"foreignKey:motorcycleId"`
	RentalPayment RentalPayment `gorm:"foreignKey:rentalId"`
	RentalRefund  RentalRefund  `gorm:"foreignKey:rentalId"`
}

func (Rental) TableName() string {
	return "rentals"
}

func (model Rental) SetReference() uint {
	return model.BaseModelUUID.ID
}
