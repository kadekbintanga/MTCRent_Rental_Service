package repository

import (
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/form"
	"service/internal/pkg/form/option"
	"service/internal/pkg/model"

	"gorm.io/gorm"
)

type RentalRepository interface {
	core.TransactionInterface
	core.FirstRepository[form.RentalFilterForm, model.Rental]
	core.FindRepository[form.RentalFilterForm, model.Rental]

	Create(form form.RentalForm, opt option.RentalOption) model.Rental
}

func NewRentalRepository(args ...*gorm.DB) RentalRepository {
	repository := rentalRepository{}
	if len(args) > 0 {
		repository.Transaction = args[0]
	}

	return &repository
}

type rentalRepository struct {
	Transaction *gorm.DB
}

func (repo *rentalRepository) SetTransaction(tx *gorm.DB) {
	repo.Transaction = tx
}

func (repo *rentalRepository) FirstByForm(form form.RentalFilterForm, args ...func(query *gorm.DB) *gorm.DB) model.Rental {
	query := repo.prepareAndFilter(form)

	if len(args) > 0 {
		query = args[0](query)
	}

	var rental model.Rental
	err := query.First(&rental).Error
	if err != nil {
		error2.ErrXtremeRentalGet(err.Error())
	}

	return rental
}

func (repo *rentalRepository) FindByForm(form form.RentalFilterForm) []model.Rental {
	query := repo.prepareAndFilter(form)

	var rental []model.Rental
	err := query.Find(&rental).Error
	if err != nil {
		error2.ErrXtremeRentalGet(err.Error())
	}

	return rental
}

func (repo *rentalRepository) Create(form form.RentalForm, opt option.RentalOption) model.Rental {
	rental := model.Rental{
		CustomerId:            opt.CustomerId,
		MotorcycleId:          opt.MotorcycleId,
		MotorcyclePlateNumber: opt.MotorcyclePlateNumber,
		RentDate:              opt.RentDate,
		ReturnDatePlan:        core.ToDate(form.ReturnDatePlan),
		PricePerDay:           opt.PricePerDay,
		StatusId:              constant.RENTAL_STATUS_ONGOING_ID,
	}

	err := repo.Transaction.Create(&rental).Error
	if err != nil {
		error2.ErrXtremeRentalSave(err.Error())
	}

	return rental
}

func (repo *rentalRepository) prepareAndFilter(form form.RentalFilterForm) *gorm.DB {
	query := config.PgSQL

	if form.ID > 0 {
		query = query.Where(`rentals."id" = ?`, form.ID)
	}

	if form.UUID != "" {
		query = query.Where(`rentals."uuid" = ?`, form.UUID)
	}

	if form.StatusId != 0 {
		query = query.Where(`rentals."statusId" = ?`, form.StatusId)
	}

	if form.CustomerId != 0 {
		query = query.Where(`rentals."customerId" = ?`, form.CustomerId)
	}

	if form.CustomerUUID != "" {
		query = query.Joins(`JOIN customers ON customers.id = rentals."customerId"`).
			Where(`customers.uuid = ?`, form.CustomerUUID)
	}

	if form.MotorcycleUUID != "" {
		query = query.Joins(`JOIN motorcycles ON motorcycles.id = rentals."motorcycleId"`).
			Where(`motorcycles.uuid = ?`, form.MotorcycleUUID)
	}

	return query
}
