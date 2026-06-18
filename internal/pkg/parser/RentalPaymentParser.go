package parser

import (
	"service/internal/pkg/constant"
	"service/internal/pkg/model"
)

type RentalPaymentParser struct {
	Array  []model.RentalPayment
	Object model.RentalPayment
}

func (parser RentalPaymentParser) Get() []interface{} {
	var result []interface{}

	for _, payment := range parser.Array {
		firstParser := RentalPaymentParser{Object: payment}
		result = append(result, firstParser.First())
	}
	return result
}

func (parser RentalPaymentParser) First() interface{} {
	payment := parser.Object

	return map[string]interface{}{
		"id":        payment.ID,
		"number":    payment.Number,
		"amount":    payment.Amount,
		"method":    constant.RentalPaymentMethod{}.IDAndName(payment.MethodId),
		"createdAt": payment.CreatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser RentalPaymentParser) CreateActivity(action string) interface{} {
	payment := parser.Object

	return map[string]interface{}{
		"id":        payment.ID,
		"number":    payment.Number,
		"rentalId":  payment.RentalId,
		"amount":    payment.Amount,
		"methodId":  constant.RentalPaymentMethod{}.IDAndName(payment.MethodId),
		"createdAt": payment.CreatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser RentalPaymentParser) UpdateActivity(action string) interface{} {
	return parser.CreateActivity(action)
}

func (parser RentalPaymentParser) DeleteActivity(action string) interface{} {
	return parser.CreateActivity(action)
}

func (parser RentalPaymentParser) GeneralActivity(action string) interface{} {
	return parser.CreateActivity(action)
}
