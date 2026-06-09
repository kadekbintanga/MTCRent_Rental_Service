package repository

import (
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/form/option"
	"service/internal/pkg/model"

	"gorm.io/gorm"
)

type RentalRefundRepository interface {
	core.TransactionInterface

	Create(opt option.RentalRefundOption) model.RentalRefund
}

func NewRentalRefundRepository(args ...*gorm.DB) RentalRefundRepository {
	repository := rentalRefundRepository{}
	if len(args) > 0 {
		repository.Transaction = args[0]
	}

	return &repository
}

type rentalRefundRepository struct {
	Transaction *gorm.DB
}

func (repo *rentalRefundRepository) SetTransaction(tx *gorm.DB) {
	repo.Transaction = tx
}

func (repo *rentalRefundRepository) Create(opt option.RentalRefundOption) model.RentalRefund {
	rentalRefund := model.RentalRefund{
		RentalId: opt.RentalId,
		Amount:   opt.Amount,
		MethodId: opt.MethodId,
	}

	err := repo.Transaction.Create(&rentalRefund).Error
	if err != nil {
		error2.ErrXtremeRentalRefundSave(err.Error())
	}

	return rentalRefund
}
