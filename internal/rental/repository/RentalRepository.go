package repository

import (
	"fmt"
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/form"
	"service/internal/pkg/form/option"
	"service/internal/pkg/model"

	"github.com/globalxtreme/go-identifier/data"
	"gorm.io/gorm"
)

type RentalRepository interface {
	core.TransactionInterface
	core.EmployeeIdentifierInterface
	core.FirstRepository[form.RentalFilterForm, model.Rental]
	core.FindRepository[form.RentalFilterForm, model.Rental]

	Create(opt option.RentalOption) model.Rental
	CountByForm(form form.RentalFilterForm) int64
	Update(rental model.Rental, opt option.RentalOption) model.Rental
	Return(rental model.Rental, form form.RentalReturnForm, opt option.RentalOption) model.Rental
}

func NewRentalRepository(args ...*gorm.DB) RentalRepository {
	repository := rentalRepository{}
	if len(args) > 0 {
		repository.tx = args[0]
	}

	return &repository
}

type rentalRepository struct {
	tx       *gorm.DB
	employee data.EmployeeIdentifierData
}

func (repo *rentalRepository) SetTransaction(tx *gorm.DB) {
	repo.tx = tx
}

func (repo *rentalRepository) SetEmployeeIdentifier(emplyee data.EmployeeIdentifierData) {
	repo.employee = emplyee
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

func (repo *rentalRepository) Create(opt option.RentalOption) model.Rental {
	rental := model.Rental{
		CustomerId:            opt.CustomerId,
		MotorcycleId:          opt.MotorcycleId,
		MotorcyclePlateNumber: opt.MotorcyclePlateNumber,
		RentDate:              opt.RentDate,
		RentDay:               opt.RentDay,
		ReturnDatePlan:        core.ToDate(opt.ReturnDatePlan),
		PricePerDay:           opt.PricePerDay,
		TotalRentPrice:        opt.TotalRentPrice,
		StatusId:              constant.RENTAL_STATUS_ONGOING_ID,
	}

	if repo.employee.ID != "" {
		rental.CreatedBy = &repo.employee.ID
		rental.CreatedByName = &repo.employee.FullName
		rental.UpdatedBy = &repo.employee.ID
		rental.UpdatedByName = &repo.employee.FullName
	}

	err := repo.tx.Create(&rental).Error
	if err != nil {
		error2.ErrXtremeRentalSave(err.Error())
	}

	return rental
}

func (repo *rentalRepository) Update(rental model.Rental, opt option.RentalOption) model.Rental {
	if opt.ReturnDatePlan != "" {
		rental.ReturnDatePlan = core.ToDate(opt.ReturnDatePlan)
		rental.RentDay = opt.RentDay
		rental.TotalRentPrice = opt.TotalRentPrice
	}

	if opt.ReturnDateActual != "" {
		rental.ReturnDateActual = core.ToDate(opt.ReturnDateActual)
		rental.LateDay = opt.LateDay
		rental.PinaltyPrice = opt.PinaltyPrice
		rental.StatusId = constant.RENTAL_STATUS_DONE_ID
	}

	rental.Note = opt.Note

	if repo.employee.ID != "" {
		rental.UpdatedBy = &repo.employee.ID
		rental.UpdatedByName = &repo.employee.FullName
	}

	err := repo.tx.Updates(&rental).Error
	if err != nil {
		error2.ErrXtremeRentalReturn(err.Error())
	}
	return rental
}

func (repo *rentalRepository) Return(rental model.Rental, form form.RentalReturnForm, opt option.RentalOption) model.Rental {
	rental.ReturnDateActual = core.ToDate(form.ReturnDateActual)
	rental.LateDay = opt.LateDay
	rental.PinaltyPrice = opt.PinaltyPrice
	rental.StatusId = constant.RENTAL_STATUS_DONE_ID
	rental.Note = form.Note

	if repo.employee.ID != "" {
		rental.UpdatedBy = &repo.employee.ID
		rental.UpdatedByName = &repo.employee.FullName
	}

	err := repo.tx.Updates(&rental).Error
	if err != nil {
		error2.ErrXtremeRentalReturn(err.Error())
	}
	return rental
}

func (repo *rentalRepository) CountByForm(form form.RentalFilterForm) int64 {
	query := repo.prepareAndFilter(form)

	var count int64
	err := query.Model(&model.Rental{}).Count(&count).Error
	if err != nil {
		error2.ErrXtremeRentalGet(err.Error())
	}
	return count
}

/** --- UNEXPORTED FUNCTIONS --- */

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

	if len(form.Orders) > 0 {
		for key, value := range form.Orders {

			query = query.Order(fmt.Sprintf("%s %s", key, value))
		}
	} else {
		query = query.Order("id DESC")
	}

	if len(form.Preloads) > 0 {
		for _, preload := range form.Preloads {
			query = query.Preload(preload)
		}
	}

	return query
}
