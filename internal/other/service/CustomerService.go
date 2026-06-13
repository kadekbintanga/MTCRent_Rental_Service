package service

import (
	"encoding/json"
	"fmt"
	"service/internal/other/repository"
	"service/internal/pkg/activity"
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	error2 "service/internal/pkg/error"
	form2 "service/internal/pkg/form"
	"service/internal/pkg/form/option"
	grpc "service/internal/pkg/grpc/customer"
	"service/internal/pkg/model"
	"service/internal/pkg/parser"
	"service/internal/pkg/saga"

	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/globalxtreme/go-identifier/data"
	"gorm.io/gorm"
)

type CustomerService interface {
	SetTransaction(tx *gorm.DB)
	SetEmployeeIdentifier(employee data.EmployeeIdentifierData)

	Save(uuid string) model.Customer
	BlacklistCustomer(customer model.Customer, reason string)
	Update(form form2.CustomerUpdateForm)
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
	srv.saga = saga.NewCustomerSaga()
	defer srv.saga.Close()

	conn := xtremepkg.RedisPool.Get()
	defer conn.Close()

	srv.repository = repository.NewCustomerRepository(srv.tx)
	customer = srv.repository.UpdateStatus(customer, option.CustomerSaveOption{StatusId: constant.CUSTOMER_STATUS_BLACKLISTED_ID})

	cacheKey := fmt.Sprintf("%s:%s", constant.CACHE_CUSTOMER, customer.UUID)
	data, _ := json.Marshal(customer)
	_, err := conn.Do("SETEX", cacheKey, constant.CACHE_TTL_COMPONENT, data)
	if err != nil {
		error2.ErrXtremeCustomerUpdate(err.Error())
	}

	srv.updateCustomerStatusSaga(customer, constant.CUSTOMER_STATUS_BLACKLISTED_ID, reason)

}

func (srv *customerService) Update(form form2.CustomerUpdateForm) {
	conn := xtremepkg.RedisPool.Get()
	defer conn.Close()

	customer := srv.prepare(&form.ID)
	cacheKey := fmt.Sprintf("%s:%s", constant.CACHE_CUSTOMER, customer.UUID)

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)
		if form.Deleted {
			if customer.ID != 0 {
				srv.repository.Delete(customer)
				_, err := conn.Do("DEL", cacheKey)
				if err != nil {
					error2.ErrXtremeCustomerDelete(err.Error())
				}
			}
		} else {
			customer = srv.repository.UpdateOrCreate(customer, form)
			data, _ := json.Marshal(customer)
			_, err := conn.Do("SETEX", cacheKey, constant.CACHE_TTL_COMPONENT, data)
			if err != nil {
				error2.ErrXtremeCustomerUpdate(err.Error())
			}
		}

		return nil
	})
}

/** --- UNEXPORTED FUNCTIONS --- */

func (srv *customerService) prepare(id *uint) model.Customer {
	srv.repository = repository.NewCustomerRepository()

	var customer model.Customer
	if id != nil {
		customer = srv.repository.FirstByForm(option.CustomerOption{
			ID: int(*id),
		})
	}

	return customer
}

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
	payload := grpc.CustomerUpdateStatusRequest{
		Uuid:            customer.UUID,
		StatusId:        int32(statusId),
		BlacklistReason: blacklistReason,
		CreatedBy:       srv.employee.ID,
		CreatedByName:   srv.employee.FullName,
	}

	srv.saga.UpdateCustomerStatus(&payload)
}
