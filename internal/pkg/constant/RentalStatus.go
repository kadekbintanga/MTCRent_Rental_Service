package constant

import "service/internal/pkg/core"

const RENTAL_STATUS_ONGOING_ID = 1
const RENTAL_STATUS_ONGOING = "Ongoing"
const RENTAL_STATUS_DONE_ID = 2
const RENTAL_STATUS_DONE = "Done"

type RentalStatus struct{}

func (in RentalStatus) OptionIDNames() map[int]string {
	return map[int]string{
		RENTAL_STATUS_ONGOING_ID: RENTAL_STATUS_ONGOING,
		RENTAL_STATUS_DONE_ID:    RENTAL_STATUS_DONE,
	}
}

func (in RentalStatus) IDAndName(id int) map[string]interface{} {
	return core.IDName{}.IDAndName(id, in)
}

func (in RentalStatus) Display(id int) string {
	return core.IDName{}.Display(id, in)
}
