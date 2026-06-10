package parser

import (
	"service/internal/pkg/constant"
	"service/internal/pkg/model"
)

type CustomerParser struct {
	Object model.Customer
}

func (parser CustomerParser) CreateActivity(action string) interface{} {
	customer := parser.Object

	return map[string]interface{}{
		"id":        customer.ID,
		"uuid":      customer.UUID,
		"name":      customer.Name,
		"idNumber":  customer.IDNumber,
		"simNumber": customer.SIMNumber,
		"phone":     customer.Phone,
		"status":    constant.CustomerStatus{}.IDAndName(customer.StatusId),
	}
}

func (parser CustomerParser) UpdateActivity(action string) interface{} {
	return parser.CreateActivity(action)
}

func (parser CustomerParser) DeleteActivity(action string) interface{} {
	return parser.CreateActivity(action)
}

func (parser CustomerParser) GeneralActivity(action string) interface{} {
	if action == constant.ACTIVITY_CUSTOMER_STATUS {
		customer := parser.Object

		return map[string]interface{}{
			"id":     customer.ID,
			"status": constant.CustomerStatus{}.IDAndName(customer.StatusId),
		}
	}

	return parser.CreateActivity(action)
}
