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
	motorcycle := srv.prepare(nil)
	if srv.checkPlateNumberDuplicate(form.PlateNumber) {
		error2.ErrXtremeMotorcycleSave("Plate Number has been registered")
	}
	if !srv.checkBrandExist(form.BrandId) {
		error2.ErrXtremeMotorcycleSave("Invalid Motorcycle Brand")
	}

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)

		motorcycle = srv.repository.Create(form)

		parser := parser.MotorcycleParser{Object: motorcycle}
		activity.UseActivity{Employee: srv.employee}.SetReference(&motorcycle).SetParser(&parser).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Create new Motorcycle: %s [%d]", motorcycle.Name, motorcycle.ID))
		return nil
	})
	return motorcycle
}

func (srv *motorcycleService) Update(uuid string, form form2.MotorcycleForm) model.Motorcycle {
	motorcycle := srv.prepare(&uuid)

	if strings.ToUpper(motorcycle.PlateNumber) != strings.ToUpper(form.PlateNumber) {
		if srv.checkPlateNumberDuplicate(form.PlateNumber) {
			error2.ErrXtremeMotorcycleSave("Plate Number has been registered")
		}
	}

	if motorcycle.BrandId != form.BrandId {
		if !srv.checkBrandExist(form.BrandId) {
			error2.ErrXtremeMotorcycleSave("Invalid Motorcycle Brand")
		}
	}

	parser := parser.MotorcycleParser{Object: motorcycle}

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)
		useActivity := activity.UseActivity{Employee: srv.employee}.SetReference(&motorcycle).SetParser(&parser).SetOldProperty(constant.ACTION_UPDATE)

		motorcycle = srv.repository.Update(motorcycle, form)

		parser.Object = motorcycle
		useActivity.SetReference(&motorcycle).SetParser(&parser).SetNewProperty(constant.ACTION_UPDATE).
			Save(fmt.Sprintf("Update Motorcycle %s [%d]", motorcycle.Name, motorcycle.ID))
		return nil
	})
	return motorcycle
}

/** --- UNEXPORTED FUNCTIONS --- */
func (srv *motorcycleService) prepare(uuid *string) model.Motorcycle {
	srv.repository = repository.NewMotorcycleRepository()

	var motorcycle model.Motorcycle
	if uuid != nil {
		motorcycle = srv.repository.FirstByForm(form2.MotorcycleFilterForm{UUID: *uuid})
	}

	return motorcycle
}

func (srv *motorcycleService) checkPlateNumberDuplicate(plateNumber string) bool {
	count := srv.repository.CountByForm(form2.MotorcycleFilterForm{PlateNumber: plateNumber})
	if count > 0 {
		return true
	}
	return false
}

func (srv *motorcycleService) checkBrandExist(id int) bool {
	brandRepo := repository.NewMotorcycleComponentBrandRepository()
	count := brandRepo.CountByForm(form2.MotorcycleComponentBrandFilterForm{ID: id})
	if count != 1 {
		return false
	}

	return true
}
