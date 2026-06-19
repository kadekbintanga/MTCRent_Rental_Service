package service

import (
	"fmt"
	"service/internal/pkg/activity"
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	error2 "service/internal/pkg/error"
	form2 "service/internal/pkg/form"
	"service/internal/pkg/form/option"
	"service/internal/pkg/model"
	"service/internal/pkg/number"
	"service/internal/pkg/parser"
	"service/internal/rental/repository"

	"github.com/globalxtreme/go-identifier/data"
	"gorm.io/gorm"
)

type RentalRefundService interface {
	SetTransaction(tx *gorm.DB)
	SetEmployeeIdentifier(employee data.EmployeeIdentifierData)

	Create(rentalUUID string, form form2.RentalRefundForm) model.RentalRefund
}

func NewRentalRefundService() RentalRefundService {
	return &rentalRefundService{}
}

type rentalRefundService struct {
	tx         *gorm.DB
	repository repository.RentalRefundRepository
	employee   data.EmployeeIdentifierData
}

func (srv *rentalRefundService) SetTransaction(tx *gorm.DB) {
	srv.tx = tx
}

func (srv *rentalRefundService) SetEmployeeIdentifier(employee data.EmployeeIdentifierData) {
	srv.employee = employee
}

func (srv *rentalRefundService) Create(rentalUUID string, form form2.RentalRefundForm) model.RentalRefund {
	rental := srv.setRental(rentalUUID)

	totalLastPayment := srv.calculatePayment(rental.Payments)
	totalLastRefund := srv.calculateRefund(rental.Refunds)

	totalRefund := totalLastRefund + form.RefundAmount
	if totalRefund > totalLastPayment {
		error2.ErrXtremeInvalidPayload("The total refund cannot be greater than the payment already made")
	}

	number := &number.RentalRedisNumber{}
	var refund model.RentalRefund
	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewRentalRefundRepository(tx)
		srv.repository.SetEmployeeIdentifier(srv.employee)
		refund = srv.repository.Create(option.RentalRefundOption{
			Number:   number.GenerateRentalRefundNumber(),
			RentalId: rental.ID,
			Amount:   form.RefundAmount,
			MethodId: form.RefundMethodId,
		})

		parser := parser.RentalRefundParser{Object: refund}
		activity.UseActivity{Employee: srv.employee}.SetReference(&refund).SetParser(&parser).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Create refund %s [%d] for rental %s [%d]", refund.Number, refund.ID, rental.Number, rental.ID))

		return nil
	})
	return refund
}

/** --- UNEXPORTED FUNCTIONS --- */
func (srv *rentalRefundService) setRental(uuid string) model.Rental {
	rentalRepository := repository.NewRentalRepository()
	rental := rentalRepository.FirstByForm(form2.RentalFilterForm{UUID: uuid, Preloads: []string{"Payments", "Refunds"}})
	if rental.StatusId != constant.RENTAL_STATUS_ONGOING_ID {
		error2.ErrXtremeInvalidPayload("Rental status is not ongoing")
	}

	return rental
}

func (srv *rentalRefundService) calculatePayment(payments []model.RentalPayment) float64 {
	var total float64
	for _, payment := range payments {
		total = total + payment.Amount
	}

	return total
}

func (srv *rentalRefundService) calculateRefund(refunds []model.RentalRefund) float64 {
	var total float64
	for _, refund := range refunds {
		total = total + refund.Amount
	}

	return total
}
