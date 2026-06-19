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

type RentalPaymentService interface {
	SetTransaction(tx *gorm.DB)
	SetEmployeeIdentifier(employee data.EmployeeIdentifierData)

	Create(rentalUUID string, form form2.RentalPaymentForm) model.RentalPayment
}

func NewRentalPaymentService() RentalPaymentService {
	return &rentalPaymentService{}
}

type rentalPaymentService struct {
	tx         *gorm.DB
	repository repository.RentalPaymentRepository
	employee   data.EmployeeIdentifierData
}

func (srv *rentalPaymentService) SetTransaction(tx *gorm.DB) {
	srv.tx = tx
}

func (srv *rentalPaymentService) SetEmployeeIdentifier(employee data.EmployeeIdentifierData) {
	srv.employee = employee
}

func (srv *rentalPaymentService) Create(rentalUUID string, form form2.RentalPaymentForm) model.RentalPayment {
	rental := srv.setRental(rentalUUID)

	number := &number.RentalRedisNumber{}
	var payment model.RentalPayment
	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewRentalPaymentRepository(tx)
		srv.repository.SetEmployeeIdentifier(srv.employee)
		payment = srv.repository.Create(option.RentalPaymentOption{
			Number:   number.GenerateRentalPaymentNumber(),
			RentalId: rental.ID,
			Amount:   form.PaymentAmount,
			MethodId: form.PaymentMethodId,
		})

		parser := parser.RentalPaymentParser{Object: payment}
		activity.UseActivity{Employee: srv.employee}.SetReference(&payment).SetParser(&parser).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Create payment %s [%d] for rental %s [%d]", payment.Number, payment.ID, rental.Number, rental.ID))

		return nil
	})
	return payment
}

/** --- UNEXPORTED FUNCTIONS --- */
func (srv *rentalPaymentService) setRental(uuid string) model.Rental {
	rentalRepository := repository.NewRentalRepository()
	rental := rentalRepository.FirstByForm(form2.RentalFilterForm{UUID: uuid, Preloads: []string{"Payments", "Refunds"}})
	if rental.StatusId != constant.RENTAL_STATUS_ONGOING_ID {
		error2.ErrXtremeInvalidPayload("Rental status is not ongoing")
	}

	return rental
}
