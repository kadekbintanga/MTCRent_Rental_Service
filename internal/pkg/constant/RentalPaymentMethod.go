package constant

import "service/internal/pkg/core"

const RENTAL_PAYMENT_METHOD_CASH_ID = 1
const RENTAL_PAYMENT_METHOD_CASH = "Cash"
const RENTAL_PAYMENT_METHOD_QRIS_ID = 2
const RENTAL_PAYMENT_METHOD_QRIS = "QRIS"
const RENTAL_PAYMENT_METHOD_TRANSFER_ID = 3
const RENTAL_PAYMENT_METHOD_TRANSFER = "Transfer"

type RentalPaymentMethod struct{}

func (in RentalPaymentMethod) OptionIDNames() map[int]string {
	return map[int]string{
		RENTAL_PAYMENT_METHOD_CASH_ID:     RENTAL_PAYMENT_METHOD_CASH,
		RENTAL_PAYMENT_METHOD_QRIS_ID:     RENTAL_PAYMENT_METHOD_QRIS,
		RENTAL_PAYMENT_METHOD_TRANSFER_ID: RENTAL_PAYMENT_METHOD_TRANSFER,
	}
}

func (in RentalPaymentMethod) IDAndName(id int) map[string]interface{} {
	return core.IDName{}.IDAndName(id, in)
}

func (in RentalPaymentMethod) Display(id int) string {
	return core.IDName{}.Display(id, in)
}
