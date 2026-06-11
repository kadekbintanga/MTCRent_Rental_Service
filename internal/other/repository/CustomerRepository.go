package repository

import (
	"gorm.io/gorm"

	"service/internal/pkg/config"
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/form"
	"service/internal/pkg/form/option"
	"service/internal/pkg/model"
)

type CustomerRepository interface {
	core.TransactionInterface
	core.FirstRepository[option.CustomerOption, model.Customer]

	Create(opt option.CustomerSaveOption) model.Customer
	UpdateStatus(customer model.Customer, opt option.CustomerSaveOption) model.Customer
	UpdateOrCreate(customer model.Customer, form form.CustomerUpdateForm) model.Customer
	Delete(customer model.Customer)
}

func NewCustomerRepository(args ...*gorm.DB) CustomerRepository {
	repository := customerRepository{}
	if len(args) > 0 {
		repository.tx = args[0]
	}

	return &repository
}

type customerRepository struct {
	tx *gorm.DB
}

func (repo *customerRepository) SetTransaction(tx *gorm.DB) {
	repo.tx = tx
}

func (repo *customerRepository) FirstByForm(opt option.CustomerOption, args ...func(query *gorm.DB) *gorm.DB) model.Customer {
	query := repo.prepareAndFilter(opt)

	if len(args) > 0 {
		query = args[0](query)
	}

	var customer model.Customer
	err := query.Find(&customer).Error
	if err != nil {
		error2.ErrXtremeCustomerGet(err.Error())
	}

	return customer
}

func (repo *customerRepository) Create(opt option.CustomerSaveOption) model.Customer {
	customer := model.Customer{
		ID:        uint(opt.ID),
		UUID:      opt.UUID,
		Name:      opt.Name,
		IDNumber:  opt.IDNumber,
		SIMNumber: opt.SIMNumber,
		Phone:     opt.Phone,
		StatusId:  opt.StatusId,
	}

	err := repo.tx.Create(&customer).Error
	if err != nil {
		error2.ErrXtremeCustomerSave(err.Error())
	}

	return customer
}

func (repo *customerRepository) UpdateStatus(customer model.Customer, opt option.CustomerSaveOption) model.Customer {
	customer.StatusId = opt.StatusId
	err := repo.tx.Updates(&customer).Error
	if err != nil {
		error2.ErrXtremeCustomerUpdate(err.Error())
	}

	return customer
}

func (repo *customerRepository) UpdateOrCreate(customer model.Customer, form form.CustomerUpdateForm) model.Customer {
	customer.ID = form.ID
	customer.UUID = form.UUID
	customer.Name = form.Name
	customer.IDNumber = form.IDNumber
	customer.SIMNumber = form.SIMNumber
	customer.Phone = form.Phone
	customer.StatusId = form.StatusId

	err := repo.tx.Save(&customer).Error
	if err != nil {
		error2.ErrXtremeCustomerUpdate(err.Error())
	}

	return customer

}

func (repo *customerRepository) Delete(customer model.Customer) {
	err := repo.tx.Delete(&customer).Error
	if err != nil {
		error2.ErrXtremeCustomerDelete(err.Error())
	}
}

func (repo *customerRepository) prepareAndFilter(opt option.CustomerOption) *gorm.DB {
	query := config.PgSQL

	if opt.ID > 0 {
		query = query.Where(`customers."id" = ?`, opt.ID)
	}

	if opt.UUID != "" {
		query = query.Where(`customers."uuid" = ?`, opt.UUID)
	}

	return query
}
