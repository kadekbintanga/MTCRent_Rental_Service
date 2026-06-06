package constant

import "service/internal/pkg/core"

const RENTAL_REFUND_METHOD_CASH_ID = 1
const RENTAL_REFUND_METHOD_CASH = "Cash"
const RENTAL_REFUND_METHOD_TRANSFER_ID = 2
const RENTAL_REFUND_METHOD_TRANSFER = "Transfer"

type RentalRefundMethod struct{}

func (in RentalRefundMethod) OptionIDNames() map[int]string {
	return map[int]string{
		RENTAL_REFUND_METHOD_CASH_ID:     RENTAL_REFUND_METHOD_CASH,
		RENTAL_REFUND_METHOD_TRANSFER_ID: RENTAL_REFUND_METHOD_TRANSFER,
	}
}

func (in RentalRefundMethod) IDAndName(id int) map[string]interface{} {
	return core.IDName{}.IDAndName(id, in)
}

func (in RentalRefundMethod) Display(id int) string {
	return core.IDName{}.Display(id, in)
}
