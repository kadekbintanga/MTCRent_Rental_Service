package service

import (
	"fmt"
	"service/internal/other/repository"
	"service/internal/pkg/activity"
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/form/option"
	grpc "service/internal/pkg/grpc/customer"
	"service/internal/pkg/model"
	"service/internal/pkg/parser"
	"service/internal/pkg/saga"

	"github.com/globalxtreme/go-identifier/data"
	"gorm.io/gorm"
)

type CustomerService interface {
	SetTransaction(tx *gorm.DB)
	SetEmployeeIdentifier(employee data.EmployeeIdentifierData)

	Save(uuid string) model.Customer
	BlacklistCustomer(customer model.Customer, reason string)
}

func NewCustomerService() CustomerService {
	return &customerService{}
}

type customerService struct {
	tx         *gorm.DB
	repository repository.CustomerRepository
	employee   data.EmployeeIdentifierData
	saga       saga.CustomerSaga
}

func (srv *customerService) SetTransaction(tx *gorm.DB) {
	srv.tx = tx
}

func (srv *customerService) SetEmployeeIdentifier(employee data.EmployeeIdentifierData) {
	srv.employee = employee
}

func (srv *customerService) Save(uuid string) model.Customer {
	dataCustomer := srv.getCustomerSaga(uuid)

	var customer model.Customer
	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCustomerRepository(tx)
		customer = srv.repository.Create(dataCustomer)

		parser := parser.CustomerParser{Object: customer}
		activity.UseActivity{Employee: srv.employee}.SetReference(&customer).SetParser(&parser).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Save new Customer [%d]", customer.ID))

		return nil
	})
	return customer

}

func (srv *customerService) BlacklistCustomer(customer model.Customer, reason string) {
	srv.repository = repository.NewCustomerRepository(srv.tx)
	srv.repository.UpdateStatus(customer, option.CustomerSaveOption{StatusId: constant.CUSTOMER_STATUS_BLACKLISTED_ID})

	srv.updateCustomerStatusSaga(customer, constant.CUSTOMER_STATUS_BLACKLISTED_ID, reason)
}

/** --- UNEXPORTED FUNCTIONS --- */

func (srv *customerService) getCustomerSaga(uuid string) option.CustomerSaveOption {
	srv.saga = saga.NewCustomerSaga()
	payload := grpc.FirstCustomerRequest{
		Uuid: uuid,
	}

	customerSaga := srv.saga.FirstCustomerByUUID(&payload)

	status := customerSaga["status"].(map[string]interface{})
	statusID, ok := status["id"].(float64)
	if !ok {
		error2.ErrXtremeCustomerSave("Invalid customer status")
	}

	if int(statusID) == constant.CUSTOMER_STATUS_BLACKLISTED_ID {
		error2.ErrXtremeCustomerSave("Customer was blacklisted")
	}

	customerID, ok := customerSaga["id"].(float64)
	if !ok {
		error2.ErrXtremeCustomerSave("Invalid customer")
	}

	return option.CustomerSaveOption{
		ID:        int(customerID),
		UUID:      customerSaga["uuid"].(string),
		Name:      customerSaga["name"].(string),
		IDNumber:  customerSaga["IDNumber"].(string),
		SIMNumber: customerSaga["SIMNumber"].(string),
		Phone:     customerSaga["phone"].(string),
		StatusId:  int(statusID),
	}
}

func (srv *customerService) updateCustomerStatusSaga(customer model.Customer, statusId int, blacklistReason string) {
	srv.saga = saga.NewCustomerSaga()
	defer srv.saga.Close()
	payload := grpc.CustomerUpdateStatusRequest{
		Uuid:            customer.UUID,
		StatusId:        int32(statusId),
		BlacklistReason: blacklistReason,
		CreatedBy:       srv.employee.ID,
		CreatedByName:   srv.employee.FullName,
	}

	srv.saga.UpdateCustomerStatus(&payload)
}
