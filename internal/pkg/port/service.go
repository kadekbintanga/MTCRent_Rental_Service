package port

import (
	"service/internal/pkg/model"

	"github.com/globalxtreme/go-identifier/data"
	"gorm.io/gorm"
)

type CustomerService interface {
	SetTransaction(tx *gorm.DB)
	SetEmployeeIdentifier(employee data.EmployeeIdentifierData)

	Save(uuid string) model.Customer
	BlacklistCustomer(customer model.Customer, reason string)
}
