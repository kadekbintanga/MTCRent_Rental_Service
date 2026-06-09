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

	return map[string]interface{}{
		"uuid":                  rental.UUID,
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

func (parser RentalParser) FirstFull() interface{} {
	rental := parser.Object

	customer := map[string]interface{}{
		"uuid":      rental.Customer.UUID,
		"name":      rental.Customer.Name,
		"idNumber":  rental.Customer.IDNumber,
		"simNumber": rental.Customer.SIMNumber,
		"phone":     rental.Customer.Phone,
	}

	motorcycle := map[string]interface{}{
		"uuid":        rental.Motorcycle.UUID,
		"plateNumber": rental.Motorcycle.PlateNumber,
		"name":        rental.Motorcycle.Name,
		"typeId":      constant.MotorcycleType{}.IDAndName(rental.Motorcycle.TypeId),
	}

	var rentalPayments []map[string]interface{}
	for _, rentalPayment := range rental.RentalPayments {
		payment := map[string]interface{}{
			"amount":   rentalPayment.Amount,
			"methodId": constant.RentalPaymentMethod{}.IDAndName(rentalPayment.MethodId),
		}
		rentalPayments = append(rentalPayments, payment)
	}

	var rentalRefunds []map[string]interface{}
	for _, rentalRefund := range rental.RentalRefunds {
		refund := map[string]interface{}{
			"amount":   rentalRefund.Amount,
			"methodId": constant.RentalRefundMethod{}.IDAndName(rentalRefund.MethodId),
		}
		rentalRefunds = append(rentalRefunds, refund)
	}

	return map[string]interface{}{
		"uuid":                  rental.UUID,
		"customer":              customer,
		"motorcycle":            motorcycle,
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
		"payment":               rentalPayments,
		"refund":                rentalRefunds,
		"createdAt":             rental.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":             rental.UpdatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser RentalParser) CreateActivity(action string) interface{} {
	rental := parser.Object

	customer := map[string]interface{}{
		"uuid":      rental.Customer.UUID,
		"name":      rental.Customer.Name,
		"idNumber":  rental.Customer.IDNumber,
		"simNumber": rental.Customer.SIMNumber,
		"phone":     rental.Customer.Phone,
	}

	motorcycle := map[string]interface{}{
		"uuid":        rental.Motorcycle.UUID,
		"plateNumber": rental.Motorcycle.PlateNumber,
		"name":        rental.Motorcycle.Name,
		"typeId":      constant.MotorcycleType{}.IDAndName(rental.Motorcycle.TypeId),
	}

	return map[string]interface{}{
		"uuid":                  rental.UUID,
		"customer":              customer,
		"motorcycle":            motorcycle,
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
