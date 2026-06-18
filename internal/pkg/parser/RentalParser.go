package parser

import (
	"service/internal/pkg/constant"
	"service/internal/pkg/model"
)

type RentalParser struct {
	Array  []model.Rental
	Object model.Rental
}

func (parser RentalParser) Get() []interface{} {
	var result []interface{}

	for _, rentals := range parser.Array {
		firstParser := RentalParser{Object: rentals}
		result = append(result, firstParser.First())
	}
	return result
}

func (parser RentalParser) First() interface{} {
	rental := parser.Object

	customerParser := CustomerParser{Object: rental.Customer}
	motorcycleParser := MotorcycleParser{Object: rental.Motorcycle}

	return map[string]interface{}{
		"uuid":                  rental.UUID,
		"number":                rental.Number,
		"customer":              customerParser.Brief(),
		"motorcycle":            motorcycleParser.Brief(),
		"motorcyclePlateNumber": rental.MotorcyclePlateNumber,
		"rentDate":              rental.RentDate.Format("02/01/2006 15:04"),
		"returnDatePlan":        rental.ReturnDatePlan.Format("02/01/2006 15:04"),
		"returnDateActual":      rental.ReturnDateActual.Format("02/01/2006 15:04"),
		"pricePerDay":           rental.PricePerDay,
		"rentDay":               rental.RentDay,
		"totalRentPrice":        rental.TotalRentPrice,
		"lateDay":               rental.LateDay,
		"pinaltyPrice":          rental.PinaltyPrice,
		"status":                constant.RentalStatus{}.IDAndName(rental.StatusId),
		"note":                  rental.Note,
		"createdBy":             rental.CreatedByName,
		"updatedBy":             rental.UpdatedByName,
		"createdAt":             rental.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":             rental.UpdatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser RentalParser) Briefs() []interface{} {
	var result []interface{}

	for _, rentals := range parser.Array {
		firstParser := RentalParser{Object: rentals}
		result = append(result, firstParser.Brief())
	}
	return result
}

func (parser RentalParser) Brief() interface{} {
	rental := parser.Object

	customerParser := CustomerParser{Object: rental.Customer}
	motorcycleParser := MotorcycleParser{Object: rental.Motorcycle}

	return map[string]interface{}{
		"uuid":                  rental.UUID,
		"number":                rental.Number,
		"customer":              customerParser.Brief(),
		"motorcycle":            motorcycleParser.Brief(),
		"motorcyclePlateNumber": rental.MotorcyclePlateNumber,
		"rentDate":              rental.RentDate.Format("02/01/2006 15:04"),
		"returnDatePlan":        rental.ReturnDatePlan.Format("02/01/2006 15:04"),
		"rentDay":               rental.RentDay,
		"totalRentPrice":        rental.TotalRentPrice,
		"pinaltyPrice":          rental.PinaltyPrice,
		"status":                constant.RentalStatus{}.IDAndName(rental.StatusId),
		"createdAt":             rental.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":             rental.UpdatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser RentalParser) CreateActivity(action string) interface{} {
	rental := parser.Object

	return map[string]interface{}{
		"uuid":                  rental.UUID,
		"number":                rental.Number,
		"motorcyclePlateNumber": rental.MotorcyclePlateNumber,
		"rentDate":              rental.RentDate.Format("02/01/2006 15:04"),
		"returnDatePlan":        rental.ReturnDatePlan.Format("02/01/2006 15:04"),
		"returnDateActual":      rental.ReturnDateActual.Format("02/01/2006 15:04"),
		"pricePerDay":           rental.PricePerDay,
		"rentDay":               rental.RentDay,
		"totalRentPrice":        rental.TotalRentPrice,
		"lateDay":               rental.LateDay,
		"pinaltyPrice":          rental.PinaltyPrice,
		"status":                constant.RentalStatus{}.IDAndName(rental.StatusId),
		"note":                  rental.Note,
		"createdAt":             rental.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":             rental.UpdatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser RentalParser) UpdateActivity(action string) interface{} {
	return parser.CreateActivity(action)
}

func (parser RentalParser) DeleteActivity(action string) interface{} {
	return parser.CreateActivity(action)
}

func (parser RentalParser) GeneralActivity(action string) interface{} {
	return parser.CreateActivity(action)
}
