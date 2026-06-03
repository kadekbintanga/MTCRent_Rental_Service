package repository

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	xtrememodel "github.com/globalxtreme/go-core/v2/model"
	"gorm.io/gorm"

	"service/internal/pkg/config"
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/form"
	"service/internal/pkg/model"
)

type MotorcycleRepository interface {
	core.TransactionInterface
	core.FirstRepository[form.MotorcycleFilterForm, model.Motorcycle]
	core.FindRepository[form.MotorcycleFilterForm, model.Motorcycle]
	core.PaginateRepository[form.MotorcycleFilterForm, model.Motorcycle]

	Create(form form.MotorcycleForm) model.Motorcycle
	Update(motorcycle model.Motorcycle, form form.MotorcycleForm) model.Motorcycle
}

func NewMotorcycleRepository(args ...*gorm.DB) MotorcycleRepository {
	repository := motorcycleRepository{}
	if len(args) > 0 {
		repository.Transaction = args[0]
	}

	return &repository
}

type motorcycleRepository struct {
	Transaction *gorm.DB
}

func (repo *motorcycleRepository) SetTransaction(tx *gorm.DB) {
	repo.Transaction = tx
}

func (repo *motorcycleRepository) FirstByForm(form form.MotorcycleFilterForm, args ...func(query *gorm.DB) *gorm.DB) model.Motorcycle {
	query := repo.prepareAndFilter(form)

	if len(args) > 0 {
		query = args[0](query)
	}

	var motorcycle model.Motorcycle
	err := query.First(&motorcycle).Error
	if err != nil {
		error2.ErrXtremeMotorcycleGet(err.Error())
	}

	return motorcycle
}

func (repo *motorcycleRepository) FindByForm(form form.MotorcycleFilterForm) []model.Motorcycle {
	query := repo.prepareAndFilter(form)

	var motorcycle []model.Motorcycle
	err := query.Find(&motorcycle).Error
	if err != nil {
		error2.ErrXtremeMotorcycleGet(err.Error())
	}

	return motorcycle
}

func (repo *motorcycleRepository) PaginateByForm(form form.MotorcycleFilterForm) ([]model.Motorcycle, interface{}) {
	parameter := url.Values{}
	parameter.Set("page", strconv.Itoa(form.Page))
	parameter.Set("limit", strconv.Itoa(form.Limit))

	query := repo.prepareAndFilter(form)
	motorcycle, pagination, err := xtrememodel.Paginate(query, parameter, model.Motorcycle{})
	if err != nil {
		error2.ErrXtremeMotorcycleGet(err.Error())
	}

	return motorcycle, pagination
}

func (repo *motorcycleRepository) Create(form form.MotorcycleForm) model.Motorcycle {
	motorcycle := model.Motorcycle{
		PlateNumber: strings.ToUpper(form.PlateNumber),
		Name:        form.Name,
		TypeId:      form.TypeId,
		Year:        form.Year,
		PricePerDay: form.PricePerDay,
		StatusId:    form.StatusId,
		BrandId:     form.BrandId,
	}

	err := repo.Transaction.Create(&motorcycle).Error
	if err != nil {
		error2.ErrXtremeMotorcycleSave(err.Error())
	}

	return motorcycle
}

func (repo *motorcycleRepository) Update(motorcycle model.Motorcycle, form form.MotorcycleForm) model.Motorcycle {
	motorcycle.PlateNumber = strings.ToUpper(form.PlateNumber)
	motorcycle.Name = form.Name
	motorcycle.TypeId = form.TypeId
	motorcycle.Year = form.Year
	motorcycle.PricePerDay = form.PricePerDay
	motorcycle.StatusId = form.StatusId
	motorcycle.BrandId = form.BrandId

	err := repo.Transaction.Updates(&motorcycle).Error
	if err != nil {
		error2.ErrXtremeMotorcycleUpdate(err.Error())
	}
	return motorcycle
}

func (repo *motorcycleRepository) UpdateStatus(motorcycle model.Motorcycle, form form.MotorcycleStatusUpdateForm) model.Motorcycle {
	err := repo.Transaction.Model(&motorcycle).Update("statusId", form.StatusId).Error
	if err != nil {
		error2.ErrXtremeMotorcycleUpdate(err.Error())
	}

	return motorcycle
}

/** --- UNEXPORTED FUNCTIONS --- */

func (repo *motorcycleRepository) prepareAndFilter(form form.MotorcycleFilterForm) *gorm.DB {
	query := config.PgSQL

	if form.ID > 0 {
		query = query.Where(`motorcycles."id" = ?`, form.ID)
	}

	if form.UUID != "" {
		query = query.Where(`motorcycles."uuid" = ?`, form.UUID)
	}

	if form.PlateNumber != "" {
		query = query.Where(`motorcycles."plateNumber"`, form.PlateNumber)
	}

	if form.TypeId != 0 {
		query = query.Where(`motorcycles."typeId" = ?`, form.TypeId)
	}

	if form.StatusId != 0 {
		query = query.Where(`motorcycles."statusId" = ?`, form.StatusId)
	}

	if search := form.Search; len(search) > 3 {
		searchVal := "%" + search + "%"
		query = query.Where(`motorcycles."name" ILIKE ? OR motorcyle."plateNumber" ILIKE ?`, searchVal, searchVal)
	}

	if len(form.Orders) > 0 {
		for key, value := range form.Orders {

			query = query.Order(fmt.Sprintf("%s %s", key, value))
		}
	} else {
		query = query.Order("id DESC")
	}

	if len(form.Preloads) > 0 {
		for _, preload := range form.Preloads {
			query = query.Preload(preload)
		}
	}

	return query
}
