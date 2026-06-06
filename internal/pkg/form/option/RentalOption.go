package option

import "time"

type RentalOption struct {
	CustomerId            int
	MotorcycleId          int
	MotorcyclePlateNumber string
	RentDate              time.Time
	PricePerDay           float64
}
