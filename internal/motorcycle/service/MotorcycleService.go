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
	"service/internal/pkg/port"
)

type MotorcycleService interface {
	SetTransaction(tx *gorm.DB)
	SetEmployeeIdentifier(employee data.EmployeeIdentifierData)
	SetActivityRepository(repo port.ActivityRepository)

	Create(form form2.MotorcycleForm) model.Motorcycle
	Update(filterForm form2.MotorcycleFilterForm, form form2.MotorcycleForm) model.Motorcycle
}

func NewMotorcycleService() MotorcycleService {
	return &motorcycleService{}
}

type motorcycleService struct {
	tx           *gorm.DB
	repository   repository.MotorcycleRepository
	activityRepo port.ActivityRepository
	employee     data.EmployeeIdentifierData
}

func (srv *motorcycleService) SetTransaction(tx *gorm.DB) {
	srv.tx = tx
}

func (srv *motorcycleService) SetEmployeeIdentifier(employee data.EmployeeIdentifierData) {
	srv.employee = employee
}

func (srv *motorcycleService) SetActivityRepository(repo port.ActivityRepository) {
	srv.activityRepo = repo
}

func (srv *motorcycleService) Create(form form2.MotorcycleForm) model.Motorcycle {
	srv.repository = repository.NewMotorcycleRepository()
	checkPlateNumber := srv.repository.FindByForm(form2.MotorcycleFilterForm{PlateNumber: form.PlateNumber})
	if len(checkPlateNumber) > 0 {
		error2.ErrXtremeMotorcycleSave("Plate Number has been registered")
	}
	brandRepo := repository.NewMotorcycleComponentBrandRepository()
	brandRepo.FirstByForm(form2.MotorcycleComponentBrandFilterForm{ID: form.BrandId})

	var motorcycle model.Motorcycle
	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewMotorcycleRepository(tx)

		motorcycle = srv.repository.Create(form)

		parser := parser.MotorcycleParser{Object: motorcycle}
		activity.UseActivity{Employee: srv.employee}.SetReference(&motorcycle).SetParser(&parser).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Create new Motorcycle: %s [%d]", motorcycle.Name, motorcycle.ID))
		return nil
	})
	return motorcycle
}

func (srv *motorcycleService) Update(filterForm form2.MotorcycleFilterForm, form form2.MotorcycleForm) model.Motorcycle {
	srv.repository = repository.NewMotorcycleRepository()
	motorcycle := srv.repository.FirstByForm(filterForm)

	if strings.ToUpper(motorcycle.PlateNumber) != strings.ToUpper(form.PlateNumber) {
		checkPlateNumber := srv.repository.FindByForm(form2.MotorcycleFilterForm{PlateNumber: form.PlateNumber})
		if len(checkPlateNumber) > 0 {
			error2.ErrXtremeMotorcycleUpdate("New Plate Number has been registered")
		}
	}

	if motorcycle.BrandId != form.BrandId {
		brandRepo := repository.NewMotorcycleComponentBrandRepository()
		brandRepo.FirstByForm(form2.MotorcycleComponentBrandFilterForm{ID: form.BrandId})
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
