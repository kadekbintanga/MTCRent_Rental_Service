package repository

import (
	"net/url"
	"service/internal/pkg/config"
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/form"
	"service/internal/pkg/form/option"
	"service/internal/pkg/model"
	"strconv"

	xtrememodel "github.com/globalxtreme/go-core/v2/model"
	"github.com/globalxtreme/go-identifier/data"
	"gorm.io/gorm"
)

type RentalPaymentRepository interface {
	core.TransactionInterface
	core.EmployeeIdentifierInterface
	core.PaginateRepository[form.RentalPaymentFilterForm, model.RentalPayment]

	Create(opt option.RentalPaymentOption) model.RentalPayment
}

func NewRentalPaymentRepository(args ...*gorm.DB) RentalPaymentRepository {
	repository := rentalPaymentRepository{}
	if len(args) > 0 {
		repository.tx = args[0]
	}

	return &repository
}

type rentalPaymentRepository struct {
	tx       *gorm.DB
	employee data.EmployeeIdentifierData
}

func (repo *rentalPaymentRepository) SetTransaction(tx *gorm.DB) {
	repo.tx = tx
}

func (repo *rentalPaymentRepository) SetEmployeeIdentifier(emplyee data.EmployeeIdentifierData) {
	repo.employee = emplyee
}

func (repo *rentalPaymentRepository) PaginateByForm(form form.RentalPaymentFilterForm) ([]model.RentalPayment, interface{}) {
	parameter := url.Values{}
	parameter.Set("page", strconv.Itoa(form.Page))
	parameter.Set("limit", strconv.Itoa(form.Limit))

	query := repo.prepareAndFilter(form)
	payments, pagination, err := xtrememodel.Paginate(query, parameter, model.RentalPayment{})
	if err != nil {
		error2.ErrXtremeRentalPaymentGet(err.Error())
	}

	return payments, pagination
}

func (repo *rentalPaymentRepository) Create(opt option.RentalPaymentOption) model.RentalPayment {
	rentalPayment := model.RentalPayment{
		Number:   opt.Number,
		RentalId: opt.RentalId,
		Amount:   opt.Amount,
		MethodId: opt.MethodId,
	}

	if repo.employee.ID != "" {
		rentalPayment.CreatedBy = &repo.employee.ID
		rentalPayment.CreatedByName = &repo.employee.FullName
	}

	err := repo.tx.Create(&rentalPayment).Error
	if err != nil {
		error2.ErrXtremeRentalPaymentSave(err.Error())
	}

	return rentalPayment
}

/** --- UNEXPORTED FUNCTIONS --- */

func (repo *rentalPaymentRepository) prepareAndFilter(form form.RentalPaymentFilterForm) *gorm.DB {
	query := config.PgSQL

	if form.ID > 0 {
		query = query.Where(`rental_payments."id" = ?`, form.ID)
	}

	if form.MethodId != 0 {
		query = query.Where(`rental_payments."methodId", ?`, form.MethodId)
	}

	if form.RentalUUID != "" {
		query = query.Joins(`JOIN rentals ON rentals.id = rental_payments."rentalId"`).
			Where(`rentals.uuid = ?`, form.RentalUUID)
	}

	if len(form.Preloads) > 0 {
		for _, preload := range form.Preloads {
			query = query.Preload(preload)
		}
	}

	query = query.Order("id DESC")

	return query
}
