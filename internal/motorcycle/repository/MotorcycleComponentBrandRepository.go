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

type MotorcycleComponentBrandRepository interface {
	core.TransactionInterface
	core.FirstRepository[form.MotorcycleComponentBrandFilterForm, model.MotorcycleComponentBrand]
	core.FindRepository[form.MotorcycleComponentBrandFilterForm, model.MotorcycleComponentBrand]
	core.PaginateRepository[form.MotorcycleComponentBrandFilterForm, model.MotorcycleComponentBrand]

	Create(form form.MotorcycleComponentBrandForm) model.MotorcycleComponentBrand
	Update(brand model.MotorcycleComponentBrand, form form.MotorcycleComponentBrandForm) model.MotorcycleComponentBrand
}

func NewMotorcycleComponentBrandRepository(args ...*gorm.DB) MotorcycleComponentBrandRepository {
	repository := motorcycleComponentBrandRepository{}
	if len(args) > 0 {
		repository.Transaction = args[0]
	}

	return &repository
}

type motorcycleComponentBrandRepository struct {
	Transaction *gorm.DB
}

func (repo *motorcycleComponentBrandRepository) SetTransaction(tx *gorm.DB) {
	repo.Transaction = tx
}

func (repo *motorcycleComponentBrandRepository) FirstByForm(form form.MotorcycleComponentBrandFilterForm, args ...func(query *gorm.DB) *gorm.DB) model.MotorcycleComponentBrand {
	query := repo.prepareAndFilter(form)

	if len(args) > 0 {
		query = args[0](query)
	}

	var motorcycleBrand model.MotorcycleComponentBrand
	err := query.First(&motorcycleBrand).Error
	if err != nil {
		error2.ErrXtremeMotorcycleBrandGet(err.Error())
	}

	return motorcycleBrand
}

func (repo *motorcycleComponentBrandRepository) FindByForm(form form.MotorcycleComponentBrandFilterForm) []model.MotorcycleComponentBrand {
	query := repo.prepareAndFilter(form)

	var motorcycleBrands []model.MotorcycleComponentBrand
	err := query.Find(&motorcycleBrands).Error
	if err != nil {
		error2.ErrXtremeMotorcycleBrandGet(err.Error())
	}

	return motorcycleBrands
}

func (repo *motorcycleComponentBrandRepository) PaginateByForm(form form.MotorcycleComponentBrandFilterForm) ([]model.MotorcycleComponentBrand, interface{}) {
	parameter := url.Values{}
	parameter.Set("page", strconv.Itoa(form.Page))
	parameter.Set("limit", strconv.Itoa(form.Limit))

	query := repo.prepareAndFilter(form)
	brand, pagination, err := xtrememodel.Paginate(query, parameter, model.MotorcycleComponentBrand{})
	if err != nil {
		error2.ErrXtremeMotorcycleBrandGet(err.Error())
	}

	return brand, pagination
}

func (repo *motorcycleComponentBrandRepository) Create(form form.MotorcycleComponentBrandForm) model.MotorcycleComponentBrand {
	motorcycleBrand := model.MotorcycleComponentBrand{
		Name: strings.ToUpper(form.Name),
	}

	err := repo.Transaction.Create(&motorcycleBrand).Error
	if err != nil {
		error2.ErrXtremeMotorcycleBrandSave(err.Error())
	}

	return motorcycleBrand
}

func (repo *motorcycleComponentBrandRepository) Update(brand model.MotorcycleComponentBrand, form form.MotorcycleComponentBrandForm) model.MotorcycleComponentBrand {
	brand.Name = strings.ToUpper(form.Name)

	err := repo.Transaction.Updates(&brand).Error
	if err != nil {
		error2.ErrXtremeMotorcycleBrandUpdate(err.Error())
	}
	return brand
}

/** --- UNEXPORTED FUNCTIONS --- */

func (repo *motorcycleComponentBrandRepository) prepareAndFilter(form form.MotorcycleComponentBrandFilterForm) *gorm.DB {
	query := config.PgSQL

	if form.ID > 0 {
		query = query.Where(`motorcycle_component_brands."id" = ?`, form.ID)
	}

	if form.Name != "" {
		query = query.Where(`UPPER(motorcycle_component_brands."name") = ?`, strings.ToUpper(form.Name))
	}

	if search := form.Search; len(search) > 3 {
		searchVal := "%" + search + "%"
		query = query.Where(`motorcycle_component_brands."name" ILIKE ?`, searchVal)
	}

	if len(form.Orders) > 0 {
		for key, value := range form.Orders {
			query = query.Order(fmt.Sprintf("%s %s", key, value))
		}
	} else {
		query = query.Order("id DESC")
	}

	return query
}
