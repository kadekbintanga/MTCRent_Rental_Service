package constant

import "service/internal/pkg/core"

const CUSTOMER_STATUS_ACTIVE_ID = 1
const CUSTOMER_STATUS_ACTIVE = "Active"
const CUSTOMER_STATUS_BLACKLISTED_ID = 2
const CUSTOMER_STATUS_BLACKLISTED = "Blacklisted"

type CustomerStatus struct{}

func (in CustomerStatus) OptionIDNames() map[int]string {
	return map[int]string{
		CUSTOMER_STATUS_ACTIVE_ID:      CUSTOMER_STATUS_ACTIVE,
		CUSTOMER_STATUS_BLACKLISTED_ID: CUSTOMER_STATUS_BLACKLISTED,
	}
}

func (in CustomerStatus) IDAndName(id int) map[string]interface{} {
	return core.IDName{}.IDAndName(id, in)
}

func (in CustomerStatus) Display(id int) string {
	return core.IDName{}.Display(id, in)
}
