package repository

import (
	"gorm.io/gorm"

	"service/internal/pkg/config"
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/form/option"
	"service/internal/pkg/model"
)

type CustomerRepository interface {
	core.TransactionInterface
	core.FirstRepository[option.CustomerOption, model.Customer]

	Create(opt option.CustomerSaveOption) model.Customer
	UpdateStatusByID(customerId uint, statusId int)
}

func NewCustomerRepository(args ...*gorm.DB) CustomerRepository {
	repository := customerRepository{}
	if len(args) > 0 {
		repository.Transaction = args[0]
	}

	return &repository
}

type customerRepository struct {
	Transaction *gorm.DB
}

func (repo *customerRepository) SetTransaction(tx *gorm.DB) {
	repo.Transaction = tx
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

	err := repo.Transaction.Create(&customer).Error
	if err != nil {
		error2.ErrXtremeCustomerSave(err.Error())
	}

	return customer
}

func (repo *customerRepository) UpdateStatusByID(customerId uint, statusId int) {
	err := repo.Transaction.Model(&model.Customer{}).Where(`customers."id" = ?`, customerId).Update("statusId", statusId).Error
	if err != nil {
		error2.ErrXtremeCustomerUpdate(err.Error())
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
