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

type MotorcycleComponentBrandService interface {
	SetTransaction(tx *gorm.DB)
	SetEmployeeIdentifier(employee data.EmployeeIdentifierData)
	SetActivityRepository(repo port.ActivityRepository)

	Create(form form2.MotorcycleComponentBrandForm) model.MotorcycleComponentBrand
	Update(filterForm form2.MotorcycleComponentBrandFilterForm, form form2.MotorcycleComponentBrandForm) model.MotorcycleComponentBrand
}

func NewMotorcycleComponentBrandService() MotorcycleComponentBrandService {
	return &motorcycleComponentBrandService{}
}

type motorcycleComponentBrandService struct {
	tx           *gorm.DB
	repository   repository.MotorcycleComponentBrandRepository
	activityRepo port.ActivityRepository
	employee     data.EmployeeIdentifierData
}

func (srv *motorcycleComponentBrandService) SetTransaction(tx *gorm.DB) {
	srv.tx = tx
}

func (srv *motorcycleComponentBrandService) SetEmployeeIdentifier(employee data.EmployeeIdentifierData) {
	srv.employee = employee
}

func (srv *motorcycleComponentBrandService) SetActivityRepository(repo port.ActivityRepository) {
	srv.activityRepo = repo
}

func (srv *motorcycleComponentBrandService) Create(form form2.MotorcycleComponentBrandForm) model.MotorcycleComponentBrand {
	srv.repository = repository.NewMotorcycleComponentBrandRepository()
	checkMotorcycleBrand := srv.repository.FindByForm(form2.MotorcycleComponentBrandFilterForm{Name: form.Name})
	if len(checkMotorcycleBrand) > 0 {
		error2.ErrXtremeMotorcycleBrandSave("Motorcycle Brand Name has been registered")
	}

	var motorcycleBrand model.MotorcycleComponentBrand
	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewMotorcycleComponentBrandRepository(tx)

		motorcycleBrand = srv.repository.Create(form)

		parser := parser.MotorcycleBrandParser{Object: motorcycleBrand}
		activity.UseActivity{Employee: srv.employee}.SetReference(&motorcycleBrand).SetParser(&parser).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Create new Motorcycle Brand: %s [%d]", motorcycleBrand.Name, motorcycleBrand.ID))
		return nil
	})
	return motorcycleBrand
}

func (srv *motorcycleComponentBrandService) Update(filterForm form2.MotorcycleComponentBrandFilterForm, form form2.MotorcycleComponentBrandForm) model.MotorcycleComponentBrand {
	srv.repository = repository.NewMotorcycleComponentBrandRepository()
	motorcycleBrand := srv.repository.FirstByForm(filterForm)

	if strings.ToUpper(motorcycleBrand.Name) != strings.ToUpper(form.Name) {
		checkMotorcycleBrandName := srv.repository.FindByForm(form2.MotorcycleComponentBrandFilterForm{Name: form.Name})
		if len(checkMotorcycleBrandName) > 0 {
			error2.ErrXtremeMotorcycleBrandUpdate("Motorcycle Brand Name has been registered")
		}
	}

	parser := parser.MotorcycleBrandParser{Object: motorcycleBrand}

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)
		useActivity := activity.UseActivity{Employee: srv.employee}.SetReference(&motorcycleBrand).SetParser(&parser).SetOldProperty(constant.ACTION_UPDATE)

		motorcycleBrand = srv.repository.Update(motorcycleBrand, form)

		parser.Object = motorcycleBrand
		useActivity.SetReference(&motorcycleBrand).SetParser(&parser).SetNewProperty(constant.ACTION_UPDATE).
			Save(fmt.Sprintf("Update Motorcycle Brand %s [%d]", motorcycleBrand.Name, motorcycleBrand.ID))
		return nil
	})
	return motorcycleBrand
}
