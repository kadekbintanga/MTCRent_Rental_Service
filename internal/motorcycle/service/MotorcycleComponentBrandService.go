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

type MotorcycleComponentBrandService interface {
	SetTransaction(tx *gorm.DB)
	SetEmployeeIdentifier(employee data.EmployeeIdentifierData)

	Create(form form2.MotorcycleComponentBrandForm) model.MotorcycleComponentBrand
	Update(id int, form form2.MotorcycleComponentBrandForm) model.MotorcycleComponentBrand
	Delete(id int)
}

func NewMotorcycleComponentBrandService() MotorcycleComponentBrandService {
	return &motorcycleComponentBrandService{}
}

type motorcycleComponentBrandService struct {
	tx         *gorm.DB
	repository repository.MotorcycleComponentBrandRepository
	employee   data.EmployeeIdentifierData
}

func (srv *motorcycleComponentBrandService) SetTransaction(tx *gorm.DB) {
	srv.tx = tx
}

func (srv *motorcycleComponentBrandService) SetEmployeeIdentifier(employee data.EmployeeIdentifierData) {
	srv.employee = employee
}

func (srv *motorcycleComponentBrandService) Create(form form2.MotorcycleComponentBrandForm) model.MotorcycleComponentBrand {
	motorcycleBrand := srv.prepareAndValidate(nil, &form, []string{})

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)

		motorcycleBrand = srv.repository.Create(form)

		parser := parser.MotorcycleBrandParser{Object: motorcycleBrand}
		activity.UseActivity{Employee: srv.employee}.SetReference(&motorcycleBrand).SetParser(&parser).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Create new Motorcycle Brand: %s [%d]", motorcycleBrand.Name, motorcycleBrand.ID))
		return nil
	})
	return motorcycleBrand
}

func (srv *motorcycleComponentBrandService) Update(id int, form form2.MotorcycleComponentBrandForm) model.MotorcycleComponentBrand {
	motorcycleBrand := srv.prepareAndValidate(&id, &form, []string{})

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

func (srv *motorcycleComponentBrandService) Delete(id int) {
	motorcycleBrand := srv.prepareAndValidate(&id, nil, []string{"Motorcycles"})

	if motorcycleBrand.Default == true {
		error2.ErrXtremeMotorcycleBrandDelete("Cannot delete default brand", nil)
	}

	if len(motorcycleBrand.Motorcycles) > 0 {
		attributes := []map[string]interface{}{}
		for _, motorcycle := range motorcycleBrand.Motorcycles {
			attributes = append(attributes, map[string]interface{}{
				"uuid":       motorcycle.UUID,
				"name":       motorcycle.Name,
				"platNumber": motorcycle.PlateNumber,
			})
		}
		error2.ErrXtremeMotorcycleBrandDelete("Cannot delete brand, because it has motorcycles", attributes)
	}

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository.SetTransaction(tx)
		srv.repository.Delete(motorcycleBrand)

		activity.UseActivity{Employee: srv.employee, Action: constant.ACTION_DELETE}.SetReference(motorcycleBrand).
			Save(fmt.Sprintf("Delete Motorcycle component Brand %s [%d]", motorcycleBrand.Name, motorcycleBrand.ID))

		return nil
	})
}

/** --- UNEXPORTED FUNCTIONS --- */

func (srv *motorcycleComponentBrandService) prepareAndValidate(id *int, form *form2.MotorcycleComponentBrandForm, preloads []string) model.MotorcycleComponentBrand {
	srv.repository = repository.NewMotorcycleComponentBrandRepository()

	var motorcycleBrand model.MotorcycleComponentBrand
	needCheckDuplicate := false

	if id != nil {
		motorcycleBrand = srv.repository.FirstByForm(
			form2.MotorcycleComponentBrandFilterForm{ID: *id, Preloads: preloads},
		)

		if form != nil && !strings.EqualFold(motorcycleBrand.Name, form.Name) {
			needCheckDuplicate = true
		}
	} else {
		if form != nil {
			needCheckDuplicate = true
		}
	}

	if needCheckDuplicate {
		if srv.repository.CountByForm(form2.MotorcycleComponentBrandFilterForm{Name: form.Name}) > 0 {
			error2.ErrXtremeInvalidPayload("Motorcycle Brand Name has been registered")
		}
	}

	return motorcycleBrand
}
