package repository

import (
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/form/option"
	"service/internal/pkg/model"

	"gorm.io/gorm"
)

type RentalPaymentRepository interface {
	core.TransactionInterface

	Create(opt option.RentalPaymentOption) model.RentalPayment
}

func NewRentalPaymentRepository(args ...*gorm.DB) RentalPaymentRepository {
	repository := rentalPaymentRepository{}
	if len(args) > 0 {
		repository.Transaction = args[0]
	}

	return &repository
}

type rentalPaymentRepository struct {
	Transaction *gorm.DB
}

func (repo *rentalPaymentRepository) SetTransaction(tx *gorm.DB) {
	repo.Transaction = tx
}

func (repo *rentalPaymentRepository) Create(opt option.RentalPaymentOption) model.RentalPayment {
	rentalPayment := model.RentalPayment{
		RentalId: opt.RentalId,
		Amount:   opt.Amount,
		MethodId: opt.MethodId,
	}

	err := repo.Transaction.Create(&rentalPayment).Error
	if err != nil {
		error2.ErrXtremeRentalPaymentSave(err.Error())
	}

	return rentalPayment
}
