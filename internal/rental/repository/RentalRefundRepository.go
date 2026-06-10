package repository

import (
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/form/option"
	"service/internal/pkg/model"

	"github.com/globalxtreme/go-identifier/data"
	"gorm.io/gorm"
)

type RentalRefundRepository interface {
	core.TransactionInterface
	core.EmployeeIdentifierInterface

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

func (repo *rentalRefundRepository) Create(opt option.RentalRefundOption) model.RentalRefund {
	rentalRefund := model.RentalRefund{
		RentalId: opt.RentalId,
		Amount:   opt.Amount,
		MethodId: opt.MethodId,
	}

	if repo.employee.ID != "" {
		rentalRefund.CreatedBy = &repo.employee.ID
		rentalRefund.CreatedByName = &repo.employee.FullName
	}

	err := repo.tx.Create(&rentalRefund).Error
	if err != nil {
		error2.ErrXtremeRentalRefundSave(err.Error())
	}

	return rentalRefund
}
