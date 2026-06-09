package parser

import (
	"service/internal/pkg/constant"
	"service/internal/pkg/model"
)

type RentalRefundParser struct {
	Array  []model.RentalRefund
	Object model.RentalRefund
}

func (parser RentalRefundParser) Get() []interface{} {
	var result []interface{}

	for _, refund := range parser.Array {
		firstParser := RentalRefundParser{Object: refund}
		result = append(result, firstParser.First())
	}
	return result
}

func (parser RentalRefundParser) First() interface{} {
	refund := parser.Object

	return map[string]interface{}{
		"id":        refund.ID,
		"amount":    refund.Amount,
		"methodId":  constant.RentalRefundMethod{}.IDAndName(refund.MethodId),
		"createdAt": refund.CreatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser RentalRefundParser) CreateActivity(action string) interface{} {
	refund := parser.Object

	return map[string]interface{}{
		"id":        refund.ID,
		"amount":    refund.Amount,
		"methodId":  constant.RentalRefundMethod{}.IDAndName(refund.MethodId),
		"createdAt": refund.CreatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser RentalRefundParser) UpdateActivity(action string) interface{} {
	return parser.CreateActivity(action)
}

func (parser RentalRefundParser) DeleteActivity(action string) interface{} {
	return parser.CreateActivity(action)
}

func (parser RentalRefundParser) GeneralActivity(action string) interface{} {
	return parser.CreateActivity(action)
}
