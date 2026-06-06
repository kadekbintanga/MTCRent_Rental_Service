package service

import (
	"service/internal/pkg/port"
	"service/internal/pkg/saga"
	"service/internal/rental/repository"

	"github.com/globalxtreme/go-identifier/data"
	"gorm.io/gorm"
)

type RentalService interface {
	SetTransaction(tx *gorm.DB)
	SetEmployeeIdentifier(employee data.EmployeeIdentifierData)
	SetCustomerRepository(repo port.CustomerRepository)
}

func NewRentalService() RentalService {
	return &rentalService{}
}

type rentalService struct {
	tx           *gorm.DB
	repository   repository.RentalRepository
	customerRepo port.CustomerRepository
	employee     data.EmployeeIdentifierData
	customerSaga saga.CustomerSaga
}

func (srv *rentalService) SetTransaction(tx *gorm.DB) {
	srv.tx = tx
}

func (srv *rentalService) SetEmployeeIdentifier(employee data.EmployeeIdentifierData) {
	srv.employee = employee
}

func (srv *rentalService) SetCustomerRepository(repo port.CustomerRepository) {
	srv.customerRepo = repo
}

// func (srv *rentalService) Create(form form2.RentalForm) model.Rental {
// 	customer := srv.customerRepo.FirstByForm(option.CustomerOption{UUID: form.CustomerUUID})
// 	if customer.ID == 0 {
// 		customerPayload := customer2.FirstCustomerRequest{
// 			Uuid: form.CustomerUUID,
// 		}
// 		srv.customerSaga = saga.NewCustomerSaga()
// 		customerSaga := srv.customerSaga.FirstCustomerByUUID(&customerPayload)
// 		if customerSaga["statusId"].(int) == constant.CUSTOMER_STATUS_BLACKLISTED_ID {
// 			error2.ErrXtremeRentalSave("Customer was blacklisted")
// 		}
// 		customer = srv.customerRepo.Create(
// 			option.CustomerSaveOption{
// 				ID:        customerSaga["id"].(int),
// 				UUID:      customerSaga["uuid"].(string),
// 				Name:      customerSaga["name"].(string),
// 				IDNumber:  customerSaga["IDNumber"].(string),
// 				SIMNumber: customerSaga["SIMNumber"].(string),
// 				Phone:     customerSaga["phone"].(string),
// 				StatusId:  customerSaga["statusId"].(int),
// 			},
// 		)

// 	} else {
// 		if customer.StatusId == constant.CUSTOMER_STATUS_BLACKLISTED_ID {
// 			error2.ErrXtremeRentalSave("Customer was blacklisted")
// 		}

// 		srv.repository = repository.NewRentalRepository()
// 		checkRental := srv.repository.FindByForm(
// 			form2.RentalFilterForm{
// 				CustomerId: int(customer.ID),
// 				StatusId:   constant.RENTAL_STATUS_ONGOING_ID,
// 			},
// 		)
// 		if len(checkRental) > 0 {
// 			error2.ErrXtremeRentalSave("Customer has ongoing rentals")
// 		}
// 	}

// }
