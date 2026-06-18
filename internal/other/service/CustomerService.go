package service

import (
	"encoding/json"
	"fmt"
	"service/internal/other/repository"
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	error2 "service/internal/pkg/error"
	form2 "service/internal/pkg/form"
	"service/internal/pkg/form/option"
	grpc "service/internal/pkg/grpc/customer"
	"service/internal/pkg/model"
	"service/internal/pkg/saga"

	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/globalxtreme/go-identifier/data"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

type CustomerService interface {
	SetTransaction(tx *gorm.DB)
	SetEmployeeIdentifier(employee data.EmployeeIdentifierData)

	FirstOrCreate(uuid string) model.Customer
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

func (srv *customerService) FirstOrCreate(uuid string) model.Customer {
	conn := xtremepkg.RedisPool.Get()
	defer conn.Close()

	var customer model.Customer
	cacheKey := fmt.Sprintf("%s:%s", constant.CACHE_CUSTOMER, uuid)
	res, err := redis.Bytes(conn.Do("GET", cacheKey))
	if err == redis.ErrNil {
		srv.repository = repository.NewCustomerRepository()
		customer = srv.repository.FirstByForm(option.CustomerOption{UUID: uuid})
		if customer.ID == 0 {
			sagaCustomer := srv.getCustomerSaga(uuid)
			config.PgSQL.Transaction(func(tx *gorm.DB) error {
				srv.repository.SetTransaction(tx)
				customer = srv.repository.Create(sagaCustomer)

				return nil
			})
		}
		if customer.StatusId == constant.CUSTOMER_STATUS_BLACKLISTED_ID {
			error2.ErrXtremeInvalidPayload("Customer was blacklisted")
		}
		if *customer.Deleted {
			error2.ErrXtremeInvalidPayload("Customer was deleted")
		}
		data, _ := json.Marshal(customer)
		_, err := conn.Do("SETEX", cacheKey, constant.CACHE_TTL_COMPONENT, data)
		if err != nil {
			error2.ErrXtremeCustomerSave("Redis : " + err.Error())
		}
	} else if err != nil {
		error2.ErrXtremeCustomerSave("Redis : " + err.Error())
	} else {
		if err := json.Unmarshal(res, &customer); err != nil {
			error2.ErrXtremeCustomerSave("Redis : " + err.Error())
		}
	}

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
	if customer.ID == 0 {
		return
	}
	cacheKey := fmt.Sprintf("%s:%s", constant.CACHE_CUSTOMER, customer.UUID)

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewCustomerRepository(tx)
		customer = srv.repository.Update(customer, form)
		if form.Deleted {
			_, err := conn.Do("DEL", cacheKey)
			if err != nil {
				error2.ErrXtremeCustomerUpdate("Redis : " + err.Error())
			}
		} else {
			data, _ := json.Marshal(customer)
			_, err := conn.Do("SETEX", cacheKey, constant.CACHE_TTL_COMPONENT, data)
			if err != nil {
				error2.ErrXtremeCustomerUpdate("Redis : " + err.Error())
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
