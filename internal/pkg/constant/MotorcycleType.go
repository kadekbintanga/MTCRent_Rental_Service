package constant

import "service/internal/pkg/core"

const MOTORCYCLE_STATUS_MATIC_ID = 1
const MOTORCYCLE_STATUS_MATIC = "Matic"
const MOTORCYCLE_STATUS_MANUAL_ID = 2
const MOTORCYCLE_STATUS_MANUAL = "Manual"
const MOTORCYCLE_STATUS_SPORT_ID = 3
const MOTORCYCLE_STATUS_SPORT = "Sport"
const MOTORCYCLE_STATUS_ELECTRIC_ID = 4
const MOTORCYCLE_STATUS_ELECTRIC = "Electric"

type MotorcycleType struct{}

func (in MotorcycleType) OptionIDNames() map[int]string {
	return map[int]string{
		MOTORCYCLE_STATUS_MATIC_ID:    MOTORCYCLE_STATUS_MATIC,
		MOTORCYCLE_STATUS_MANUAL_ID:   MOTORCYCLE_STATUS_MANUAL,
		MOTORCYCLE_STATUS_SPORT_ID:    MOTORCYCLE_STATUS_SPORT,
		MOTORCYCLE_STATUS_ELECTRIC_ID: MOTORCYCLE_STATUS_ELECTRIC,
	}
}

func (in MotorcycleType) IDAndName(id int) map[string]interface{} {
	return core.IDName{}.IDAndName(id, in)
}

func (in MotorcycleType) Display(id int) string {
	return core.IDName{}.Display(id, in)
}
