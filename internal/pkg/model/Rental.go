package model

import (
	"time"

	xtrememodel "github.com/globalxtreme/go-core/v2/model"
)

type Rental struct {
	xtrememodel.BaseModelUUID
	CustomerId            uint      `gorm:"column:customerId;type:int;not null"`
	MotorcycleId          uint      `gorm:"column:motorcycleId;type:int;not null"`
	MotorcyclePlateNumber string    `gorm:"column:motorcyclePlateNumber;type:varchar(250);not null"`
	RentDate              time.Time `gorm:"column:rentDate;type:date;not null"`
	ReturnDatePlan        time.Time `gorm:"column:returnDatePlan;type:date;not null"`
	ReturnDateActual      time.Time `gorm:"column:returnDateActual;type:date;default:null"`
	PricePerDay           float64   `gorm:"column:pricePerDay;type:float;not null"`
	RentDay               uint      `gorm:"column:rentDay;type:int"`
	TotalRentPrice        float64   `gorm:"column:totalRentPrice;type:float"`
	LateDay               uint      `gorm:"column:lateDay;type:int"`
	PinaltyPrice          float64   `gorm:"column:pinaltyPrice;type:float"`
	StatusId              int       `gorm:"column:statusId;type:int;not null"`
	Note                  string    `gorm:"column:note;type:text;default:null"`
	CreatedBy             *string   `gorm:"column:createdBy;varchar(50);null"`
	CreatedByName         *string   `gorm:"column:createdByName;varchar(250);null"`
	UpdatedBy             *string   `gorm:"column:updatedBy;varchar(50);null"`
	UpdatedByName         *string   `gorm:"column:updatedByName;varchar(250);null"`

	Customer       Customer        `gorm:"foreignKey:CustomerId"`
	Motorcycle     Motorcycle      `gorm:"foreignKey:MotorcycleId"`
	RentalPayments []RentalPayment `gorm:"foreignKey:RentalId;references:ID"`
	RentalRefunds  []RentalRefund  `gorm:"foreignKey:RentalId;references:ID"`
}

func (Rental) TableName() string {
	return "rentals"
}

func (model Rental) SetReference() uint {
	return model.BaseModelUUID.ID
}
