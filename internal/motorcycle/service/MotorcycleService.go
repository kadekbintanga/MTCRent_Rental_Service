package service

import (
	"fmt"
	"strings"

	"github.com/globalxtreme/go-identifier/data"
	"gorm.io/gorm"

	"service/internal/motorcycle/repository"
	"service/internal/pkg/activity"
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	error2 "service/internal/pkg/error"
	form2 "service/internal/pkg/form"
	"service/internal/pkg/model"
	"service/internal/pkg/parser"
)

type MotorcycleService interface {
	SetTransaction(tx *gorm.DB)
	SetEmployeeIdentifier(employee data.EmployeeIdentifierData)

	Create(form form2.MotorcycleForm) model.Motorcycle
	Update(uuid string, form form2.MotorcycleForm) model.Motorcycle
	Delete(uuid string)
}

func NewMotorcycleService() MotorcycleService {
	return &motorcycleService{}
}

type motorcycleService struct {
	tx         *gorm.DB
	repository repository.MotorcycleRepository
	employee   data.EmployeeIdentifierData
}

func (srv *motorcycleService) SetTransaction(tx *gorm.DB) {
	srv.tx = tx
}

func (srv *motorcycleService) SetEmployeeIdentifier(employee data.EmployeeIdentifierData) {
	srv.employee = employee
}

func (srv *motorcycleService) Create(form form2.MotorcycleForm) model.Motorcycle {
	motorcycle, brand := srv.prepareAndValidate(nil, &form, []string{})

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)

		motorcycle = srv.repository.Create(form)
		motorcycle.Brand = brand

		parser := parser.MotorcycleParser{Object: motorcycle}
		activity.UseActivity{Employee: srv.employee}.SetReference(&motorcycle).SetParser(&parser).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Create new Motorcycle: %s [%d]", motorcycle.Name, motorcycle.ID))
		return nil
	})
	return motorcycle
}

func (srv *motorcycleService) Update(uuid string, form form2.MotorcycleForm) model.Motorcycle {
	motorcycle, brand := srv.prepareAndValidate(&uuid, &form, []string{})

	parser := parser.MotorcycleParser{Object: motorcycle}

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)
		useActivity := activity.UseActivity{Employee: srv.employee}.SetReference(&motorcycle).SetParser(&parser).SetOldProperty(constant.ACTION_UPDATE)

		motorcycle = srv.repository.Update(motorcycle, form)
		motorcycle.Brand = brand

		parser.Object = motorcycle
		useActivity.SetReference(&motorcycle).SetParser(&parser).SetNewProperty(constant.ACTION_UPDATE).
			Save(fmt.Sprintf("Update Motorcycle %s [%d]", motorcycle.Name, motorcycle.ID))
		return nil
	})
	return motorcycle
}

func (srv *motorcycleService) Delete(uuid string) {
	motorcycle, _ := srv.prepareAndValidate(&uuid, nil, []string{"Rentals"})
	if len(motorcycle.Rentals) > 0 {
		error2.ErrXtremeInvalidPayload("Cannot delete motorcycle that has been rented")
	}

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)

		srv.repository.Delete(motorcycle)

		activity.UseActivity{Employee: srv.employee, Action: constant.ACTION_DELETE}.SetReference(motorcycle).
			Save(fmt.Sprintf("Delete Motorcycle %s [%d]", motorcycle.Name, motorcycle.ID))

		return nil
	})
}

/** --- UNEXPORTED FUNCTIONS --- */
func (srv *motorcycleService) prepareAndValidate(uuid *string, form *form2.MotorcycleForm, preloads []string) (model.Motorcycle, model.MotorcycleComponentBrand) {
	srv.repository = repository.NewMotorcycleRepository()
	srv.repository.SetEmployeeIdentifier(srv.employee)

	var motorcycle model.Motorcycle
	var motorcycleBrand model.MotorcycleComponentBrand
	needCheckPlateNumber := false
	needCheckBrand := false

	if uuid != nil {
		motorcycle = srv.repository.FirstByForm(form2.MotorcycleFilterForm{UUID: *uuid, Preloads: preloads})
		if form != nil {
			if !strings.EqualFold(motorcycle.PlateNumber, form.PlateNumber) {
				needCheckPlateNumber = true
			}
			if motorcycle.BrandId != form.BrandId {
				needCheckBrand = true
			}
		}
	} else {
		if form != nil {
			needCheckPlateNumber = true
			needCheckBrand = true
		}
	}

	if needCheckPlateNumber {
		count := srv.repository.CountByForm(form2.MotorcycleFilterForm{PlateNumber: form.PlateNumber})
		if count > 0 {
			error2.ErrXtremeInvalidPayload("Plate Number has been registered")
		}
	}

	if needCheckBrand {
		brandRepo := repository.NewMotorcycleComponentBrandRepository()
		motorcycleBrand = brandRepo.FirstByForm(form2.MotorcycleComponentBrandFilterForm{ID: form.BrandId})
	}

	return motorcycle, motorcycleBrand
}
