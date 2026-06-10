package repository

import (
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/form/option"
	"service/internal/pkg/model"

	"github.com/globalxtreme/go-identifier/data"
	"gorm.io/gorm"
)

type RentalPaymentRepository interface {
	core.TransactionInterface
	core.EmployeeIdentifierInterface

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

func (repo *rentalPaymentRepository) Create(opt option.RentalPaymentOption) model.RentalPayment {
	rentalPayment := model.RentalPayment{
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
