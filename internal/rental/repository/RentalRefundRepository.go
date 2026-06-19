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

type RentalRefundRepository interface {
	core.TransactionInterface
	core.EmployeeIdentifierInterface
	core.PaginateRepository[form.RentalRefundFilterForm, model.RentalRefund]

	Create(opt option.RentalRefundOption) model.RentalRefund
}

func NewRentalRefundRepository(args ...*gorm.DB) RentalRefundRepository {
	repository := rentalRefundRepository{}
	if len(args) > 0 {
		repository.tx = args[0]
	}

	return &repository
}

type rentalRefundRepository struct {
	tx       *gorm.DB
	employee data.EmployeeIdentifierData
}

func (repo *rentalRefundRepository) SetTransaction(tx *gorm.DB) {
	repo.tx = tx
}

func (repo *rentalRefundRepository) SetEmployeeIdentifier(emplyee data.EmployeeIdentifierData) {
	repo.employee = emplyee
}

func (repo *rentalRefundRepository) PaginateByForm(form form.RentalRefundFilterForm) ([]model.RentalRefund, interface{}) {
	parameter := url.Values{}
	parameter.Set("page", strconv.Itoa(form.Page))
	parameter.Set("limit", strconv.Itoa(form.Limit))

	query := repo.prepareAndFilter(form)
	refunds, pagination, err := xtrememodel.Paginate(query, parameter, model.RentalRefund{})
	if err != nil {
		error2.ErrXtremeRentalRefundGet(err.Error())
	}

	return refunds, pagination
}

func (repo *rentalRefundRepository) Create(opt option.RentalRefundOption) model.RentalRefund {
	rentalRefund := model.RentalRefund{
		Number:        opt.Number,
		RentalId:      opt.RentalId,
		Amount:        opt.Amount,
		MethodId:      opt.MethodId,
		CreatedBy:     &repo.employee.ID,
		CreatedByName: &repo.employee.FullName,
	}

	err := repo.tx.Create(&rentalRefund).Error
	if err != nil {
		error2.ErrXtremeRentalRefundSave(err.Error())
	}

	return rentalRefund
}

/** --- UNEXPORTED FUNCTIONS --- */

func (repo *rentalRefundRepository) prepareAndFilter(form form.RentalRefundFilterForm) *gorm.DB {
	query := config.PgSQL

	if form.ID > 0 {
		query = query.Where(`rental_refunds."id" = ?`, form.ID)
	}

	if form.MethodId != 0 {
		query = query.Where(`rental_refunds."methodId", ?`, form.MethodId)
	}

	if form.RentalUUID != "" {
		query = query.Joins(`JOIN rentals ON rentals.id = rental_refunds."rentalId"`).
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
