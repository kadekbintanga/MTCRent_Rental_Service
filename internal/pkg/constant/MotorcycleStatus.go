package constant

import "service/internal/pkg/core"

const MOTORCYCLE_STATUS_AVAILABLE_ID = 1
const MOTORCYCLE_STATUS_AVAILABLE = "Available"
const MOTORCYCLE_STATUS_RENTED_ID = 2
const MOTORCYCLE_STATUS_RENTED = "Rented"
const MOTORCYCLE_STATUS_UNAVAILABLE_ID = 3
const MOTORCYCLE_STATUS_UNAVAILABLE = "Unavailable"
const MOTORCYCLE_STATUS_DISCONTINUED_ID = 4
const MOTORCYCLE_STATUS_DISCONTINUED = "Discontinued"

type MotorcycleStatus struct{}

func (in MotorcycleStatus) OptionIDNames() map[int]string {
	return map[int]string{
		MOTORCYCLE_STATUS_AVAILABLE_ID:    MOTORCYCLE_STATUS_AVAILABLE,
		MOTORCYCLE_STATUS_RENTED_ID:       MOTORCYCLE_STATUS_RENTED,
		MOTORCYCLE_STATUS_UNAVAILABLE_ID:  MOTORCYCLE_STATUS_UNAVAILABLE,
		MOTORCYCLE_STATUS_DISCONTINUED_ID: MOTORCYCLE_STATUS_DISCONTINUED,
	}
}

func (in MotorcycleStatus) IDAndName(id int) map[string]interface{} {
	return core.IDName{}.IDAndName(id, in)
}

func (in MotorcycleStatus) Display(id int) string {
	return core.IDName{}.Display(id, in)
}
