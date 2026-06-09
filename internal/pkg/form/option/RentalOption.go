package option

import "time"

type RentalOption struct {
	CustomerId            uint
	MotorcycleId          uint
	MotorcyclePlateNumber string
	RentDate              time.Time
	RentDay               uint
	PricePerDay           float64
	TotalRentPrice        float64
	LateDay               uint
	PinaltyPrice          float64
	ReturnDatePlan        string
	ReturnDateActual      string
	Note                  string
}
