package service

import (
	"fmt"
	"math"
	"service/internal/pkg/activity"
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
	form2 "service/internal/pkg/form"
	"service/internal/pkg/form/option"
	"service/internal/pkg/model"
	"service/internal/pkg/parser"
	"service/internal/pkg/port"
	"service/internal/pkg/saga"
	"service/internal/rental/repository"
	"time"

	xtremerabbitmq "github.com/globalxtreme/go-core/v2/rabbitmq"
	"github.com/globalxtreme/go-identifier/data"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
)

type RentalService interface {
	SetTransaction(tx *gorm.DB)
	SetEmployeeIdentifier(employee data.EmployeeIdentifierData)
	SetCustomerRepository(repo port.CustomerRepository)
	SetMotorcycleRepository(repo port.MotorcycleRepository)
	SetSettingConfigurationRepository(repo port.SettingConfigurationRepository)
	SetCustomerService(service port.CustomerService)

	Create(form form2.RentalForm) model.Rental
	Simulate(uuid string, form form2.RentalSimulateForm) map[string]interface{}
	Update(uuid string, form form2.RentalUpdateteForm) model.Rental
	Refund(uuid string, form form2.RentalRefundForm) model.RentalRefund
	Return(uuid string, form form2.RentalReturnForm) model.Rental
}

func NewRentalService() RentalService {
	return &rentalService{}
}

type rentalService struct {
	tx                *gorm.DB
	repository        repository.RentalRepository
	customerRepo      port.CustomerRepository
	motorcycleRepo    port.MotorcycleRepository
	settingConfigRepo port.SettingConfigurationRepository
	employee          data.EmployeeIdentifierData
	customerSaga      saga.CustomerSaga
	customerService   port.CustomerService
}

func (srv *rentalService) SetTransaction(tx *gorm.DB) {
	srv.tx = tx
}

func (srv *rentalService) SetEmployeeIdentifier(employee data.EmployeeIdentifierData) {
	srv.employee = employee
}

func (srv *rentalService) SetCustomerRepository(repo port.CustomerRepository) {
	srv.customerRepo = repo
}

func (srv *rentalService) SetMotorcycleRepository(repo port.MotorcycleRepository) {
	srv.motorcycleRepo = repo
}

func (srv *rentalService) SetSettingConfigurationRepository(repo port.SettingConfigurationRepository) {
	srv.settingConfigRepo = repo
}

func (srv *rentalService) SetCustomerService(service port.CustomerService) {
	srv.customerService = service
}

func (srv *rentalService) Create(form form2.RentalForm) model.Rental {
	customer := srv.checkCustomer(form.CustomerUUID)

	rental := srv.prepare(nil, []string{})

	srv.checkCustomerHasOngoingRental(form.CustomerUUID)

	motorcycle := srv.checkMotorcycle(form.MotorcycleUUID)

	rentDate := time.Now()
	rentDay, totalRentPrice := srv.calculateTotalRentPrice(rentDate, form.ReturnDatePlan, motorcycle.PricePerDay)

	rentalOpt := option.RentalOption{
		CustomerId:            customer.ID,
		MotorcycleId:          motorcycle.ID,
		MotorcyclePlateNumber: motorcycle.PlateNumber,
		RentDate:              rentDate,
		RentDay:               uint(rentDay),
		PricePerDay:           motorcycle.PricePerDay,
		TotalRentPrice:        totalRentPrice,
		ReturnDatePlan:        form.ReturnDatePlan,
	}

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)
		paymentRepository := repository.NewRentalPaymentRepository(tx)
		srv.motorcycleRepo.SetTransaction(tx)

		rental = srv.repository.Create(rentalOpt)

		payment := paymentRepository.Create(option.RentalPaymentOption{
			RentalId: rental.ID,
			Amount:   form.PaymentAmount,
			MethodId: form.PaymentMethodId,
		})

		rental.Customer = customer
		rental.RentalPayments = append(rental.RentalPayments, payment)

		motorcycle := srv.motorcycleRepo.UpdateStatus(motorcycle, form2.MotorcycleStatusUpdateForm{StatusId: constant.MOTORCYCLE_STATUS_RENTED_ID})

		rental.Motorcycle = motorcycle

		parser := parser.RentalParser{Object: rental}
		activity.UseActivity{Employee: srv.employee}.SetReference(&rental).SetParser(&parser).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Create new Rental [%d]", rental.ID))

		return nil
	})
	return rental

}

func (srv *rentalService) Simulate(uuid string, form form2.RentalSimulateForm) map[string]interface{} {
	rental := srv.prepare(&uuid, []string{"RentalPayments", "RentalRefunds"})
	if rental.StatusId != constant.RENTAL_STATUS_ONGOING_ID {
		error2.ErrXtremeRentalSimulate("Rental status is not ongoing")
	}

	lateDay, pinaltyPrice, totalPrice := srv.calculateTotalPrice(rental, form.ReturnDate)

	totalLastPayment := srv.calculatePayment(rental.RentalPayments)
	totalLastRefund := srv.calculateRefund(rental.RentalRefunds)

	var needPayment float64
	var needRefund float64

	remainingPayment := totalPrice - (totalLastPayment - totalLastRefund)

	switch {
	case remainingPayment > 0:
		needPayment = remainingPayment
	case remainingPayment < 0:
		needRefund = math.Abs(remainingPayment)
	default:
		needPayment = 0
		needRefund = 0
	}

	return map[string]interface{}{
		"totalRentPrice":   rental.TotalRentPrice,
		"lateDay":          lateDay,
		"pinaltyPrice":     pinaltyPrice,
		"totalPrice":       totalPrice,
		"totalLastPayment": totalLastPayment,
		"totalLastRefund":  totalLastRefund,
		"needPayment":      needPayment,
		"needRefund":       needRefund,
	}
}

func (srv *rentalService) Update(uuid string, form form2.RentalUpdateteForm) model.Rental {
	rental := srv.prepare(&uuid, []string{})
	if rental.StatusId != constant.RENTAL_STATUS_ONGOING_ID {
		error2.ErrXtremeRentalSimulate("Rental status is not ongoing")
	}

	rentDay, totalRentPrice := srv.calculateTotalRentPrice(rental.RentDate, form.ReturnDatePlan, rental.PricePerDay)

	rentalOpt := option.RentalOption{
		RentDay:        uint(rentDay),
		TotalRentPrice: totalRentPrice,
		ReturnDatePlan: form.ReturnDatePlan,
	}

	parser := parser.RentalParser{Object: rental}

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)
		useActivity := activity.UseActivity{Employee: srv.employee}.SetReference(&rental).SetParser(&parser).SetOldProperty(constant.ACTION_UPDATE)

		rental = srv.repository.Update(rental, rentalOpt)

		parser.Object = rental
		useActivity.SetReference(&rental).SetParser(&parser).SetNewProperty(constant.ACTION_UPDATE).
			Save(fmt.Sprintf("Update rental [%d]", rental.ID))
		return nil
	})
	return rental
}

func (srv *rentalService) Refund(uuid string, form form2.RentalRefundForm) model.RentalRefund {
	rental := srv.prepare(&uuid, []string{"RentalPayments", "RentalRefunds"})
	if rental.StatusId != constant.RENTAL_STATUS_ONGOING_ID {
		error2.ErrXtremeRentalSimulate("Rental status is not ongoing")
	}

	totalLastPayment := srv.calculatePayment(rental.RentalPayments)
	totalLastRefund := srv.calculateRefund(rental.RentalRefunds)

	totalRefund := totalLastRefund + form.RefundAmount

	if totalRefund > totalLastPayment {
		error2.ErrXtremeRentalRefund("The total refund cannot be greater than the payment already made")
	}
	var refund model.RentalRefund
	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		refundRepository := repository.NewRentalRefundRepository(tx)

		refund = refundRepository.Create(option.RentalRefundOption{
			RentalId: rental.ID,
			Amount:   form.RefundAmount,
			MethodId: form.RefundMethodId,
		})

		parser := parser.RentalRefundParser{Object: refund}
		activity.UseActivity{Employee: srv.employee}.SetReference(&refund).SetParser(&parser).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Create refund [%d] for rental [%d]", refund.ID, rental.ID))

		return nil
	})

	return refund

}

func (srv *rentalService) Return(uuid string, form form2.RentalReturnForm) model.Rental {
	rental := srv.prepare(&uuid, []string{"Motorcycle", "Customer", "RentalPayments", "RentalRefunds"})
	if rental.StatusId != constant.RENTAL_STATUS_ONGOING_ID {
		error2.ErrXtremeRentalReturn("Rental status is not ongoing")
	}

	lateDay, pinaltyPrice, totalPrice := srv.calculateTotalPrice(rental, form.ReturnDateActual)

	totalLastPayment := srv.calculatePayment(rental.RentalPayments)
	totalLastRefund := srv.calculateRefund(rental.RentalRefunds)

	remainingPayment := totalPrice - (totalLastPayment - totalLastRefund)
	if remainingPayment < 0 {
		error2.ErrXtremeRentalReturn(fmt.Sprintf("Need refund %.0f", math.Abs(remainingPayment)))
	}

	if remainingPayment != form.PaymentAmount {
		error2.ErrXtremeRentalReturn(fmt.Sprintf("Payment amount should be %.0f", remainingPayment))
	}

	rentalOpt := option.RentalOption{
		LateDay:          uint(lateDay),
		PinaltyPrice:     pinaltyPrice,
		ReturnDateActual: form.ReturnDateActual,
		Note:             form.Note,
	}

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)
		srv.motorcycleRepo.SetTransaction(tx)

		rental = srv.repository.Update(rental, rentalOpt)

		motorcycle := srv.motorcycleRepo.UpdateStatus(rental.Motorcycle, form2.MotorcycleStatusUpdateForm{StatusId: constant.MOTORCYCLE_STATUS_AVAILABLE_ID})
		if form.PaymentAmount != 0 {
			paymentRepository := repository.NewRentalPaymentRepository(tx)
			payment := paymentRepository.Create(option.RentalPaymentOption{
				RentalId: rental.ID,
				Amount:   form.PaymentAmount,
				MethodId: form.PaymentMethodId,
			})
			rental.RentalPayments = append(rental.RentalPayments, payment)
		}

		rental.Motorcycle = motorcycle
		srv.customerService.SetTransaction(tx)
		srv.customerService.SetEmployeeIdentifier(srv.employee)
		srv.processBlacklist(rental.Customer, lateDay)

		parser := parser.RentalParser{Object: rental}

		activity.UseActivity{Employee: srv.employee}.SetReference(&rental).SetParser(&parser).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Return Rental [%d]", rental.ID))

		return nil
	})
	return rental
}

/** --- UNEXPORTED FUNCTIONS --- */

func (srv *rentalService) prepare(uuid *string, preloads []string) model.Rental {
	srv.repository = repository.NewRentalRepository()
	srv.repository.SetEmployeeIdentifier(srv.employee)

	var rental model.Rental
	if uuid != nil {
		rental = srv.repository.FirstByForm(form2.RentalFilterForm{UUID: *uuid, Preloads: preloads})
	}
	return rental
}

func (srv *rentalService) checkCustomer(customerUUID string) model.Customer {
	customer := srv.customerRepo.FirstByForm(option.CustomerOption{UUID: customerUUID})
	if customer.ID == 0 {
		customer = srv.customerService.Save(customerUUID)

	} else {
		if customer.StatusId == constant.CUSTOMER_STATUS_BLACKLISTED_ID {
			error2.ErrXtremeRentalSave("Customer was blacklisted")
		}
	}

	return customer
}

func (srv *rentalService) checkCustomerHasOngoingRental(customerUUID string) {
	count := srv.repository.CountByForm(form2.RentalFilterForm{CustomerUUID: customerUUID, StatusId: constant.RENTAL_STATUS_ONGOING_ID})
	if count > 0 {
		error2.ErrXtremeRentalSave("Customer has ongoing rental")
	}
}

func (srv *rentalService) checkMotorcycle(motorcycleUUID string) model.Motorcycle {
	motorcycle := srv.motorcycleRepo.FirstByForm(form2.MotorcycleFilterForm{UUID: motorcycleUUID})
	if motorcycle.StatusId != constant.MOTORCYCLE_STATUS_AVAILABLE_ID {
		error2.ErrXtremeMotorcycleGet("Motorcycle is not available")
	}
	return motorcycle
}

func (srv *rentalService) calculatePinalty(pricePerDay float64, lateDay int) float64 {
	pinaltyPerDay := (pricePerDay * 50 / 100) + pricePerDay
	return pinaltyPerDay * float64(lateDay)
}

func (srv *rentalService) calculatePayment(payments []model.RentalPayment) float64 {
	var total float64
	for _, payment := range payments {
		total = total + payment.Amount
	}

	return total
}

func (srv *rentalService) calculateRefund(refunds []model.RentalRefund) float64 {
	var total float64
	for _, refund := range refunds {
		total = total + refund.Amount
	}

	return total
}

func (srv *rentalService) calculateTotalPrice(rental model.Rental, returnDate string) (int, float64, float64) {
	timeReturn := core.ToDate(returnDate)
	if timeReturn.Before(rental.ReturnDatePlan) {
		error2.ErrXtremeInvalidPayload("Return Date cannot be before return date plan")

	}
	lateDay, err := core.DaysUntil(rental.ReturnDatePlan, returnDate)
	if err != nil {
		error2.ErrXtremeInvalidPayload(err.Error())
	}
	var pinaltyPrice float64
	totalPrice := rental.TotalRentPrice
	if lateDay > 0 {
		pinaltyPrice = srv.calculatePinalty(rental.PricePerDay, lateDay)
		totalPrice = totalPrice + pinaltyPrice
	}
	return lateDay, pinaltyPrice, totalPrice

}

func (srv *rentalService) processBlacklist(customer model.Customer, lateDay int) {
	blacklistLimitDay := srv.settingConfigRepo.FirstByForm(form2.SettingConfigurationFilterForm{Key: constant.SETTING_CONFIGURATION_KEY_BLACKLIST_LIMIT_DAY})
	limitDay := core.ToInt(blacklistLimitDay.Value)
	if lateDay > limitDay {
		srv.customerService.BlacklistCustomer(customer, constant.CUSTOMER_PASS_DAY_LIMIT_REASON)
	}
}

func (srv *rentalService) sendBlacklistCustomer(customer model.Customer) {
	async := xtremerabbitmq.GXAsyncWorkflow{
		Action:        constant.ASYNC_WORKFLOW_CUSTOMER_STATUS_BLACKLISTED_ACTION,
		Description:   fmt.Sprintf("Blacklist customer for id : %d", customer.ID),
		ReferenceId:   utils.ToString(customer.ID),
		ReferenceType: customer.TableName(),
		Strict:        false,
	}

	async.OnStep(xtremerabbitmq.GXAsyncWorkflowStepOpt{
		Service:     constant.ASYNC_WORKFLOW_SERVICE_CUSTOMER,
		Queue:       constant.RABBITMQ_EXCHANGE_RENTAL_CUSTOMER_STATUS_UPDATE,
		Description: "Update customer status in customer service",
		Payload: map[string]interface{}{
			"uuid":            customer.UUID,
			"statusId":        constant.CUSTOMER_STATUS_BLACKLISTED_ID,
			"blacklistReason": constant.CUSTOMER_PASS_DAY_LIMIT_REASON,
		},
	})

	err := core.ErrorAsyncHandler(func() error {
		async.Push()
		return nil
	})

	if err != nil {
		error2.ErrXtremeAsyncWorkflowPush(err.Error())
	}
}

func (srv *rentalService) calculateTotalRentPrice(rentDate time.Time, returnDate string, pricePerDay float64) (int, float64) {
	rentDay, err := core.DaysUntil(rentDate, returnDate)
	if err != nil {
		error2.ErrXtremeRentalSave(err.Error())
	}

	if rentDay < 0 {
		error2.ErrXtremeRentalSave("Invalid return date plan")
	}

	totalRentPrice := float64(rentDay) * pricePerDay

	return rentDay, totalRentPrice
}
