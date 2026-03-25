package repository

import (
	"errors"
	"fmt"
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
	"gorm.io/gorm"
	"net/url"
	"service/internal/pkg/config"
	"service/internal/pkg/core"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/form"
	"service/internal/pkg/model"
	"strconv"
)

/** --- INTERFACE --- */

// TODO: Hanya contoh. nanti langsung hapus saja
type TestingRepository interface {
	core.TransactionInterface
	core.FirstRepository[form.TestingFilterForm, model.Testing]
	core.FindRepository[form.TestingFilterForm, model.Testing]
	core.PaginateRepository[form.TestingFilterForm, model.Testing]

	Create(form form.TestingForm) model.Testing
	Delete(testing model.Testing)

	AddSub(testing model.Testing, sub string) model.TestingSub
	DeleteSub(testingSub model.TestingSub)
}

func NewTestingRepository(args ...*gorm.DB) TestingRepository {
	repository := testingRepository{}
	if len(args) > 0 {
		repository.transaction = args[0]
	}

	return &repository
}

/** --- MAIN REPOSITORY --- */

type testingRepository struct {
	transaction *gorm.DB
}

func (repo *testingRepository) SetTransaction(tx *gorm.DB) {
	repo.transaction = tx
}

func (repo *testingRepository) FirstByForm(form form.TestingFilterForm, args ...func(query *gorm.DB) *gorm.DB) model.Testing {
	var testing model.Testing

	query := config.PgSQL
	if len(args) > 0 {
		query = args[0](query)
	}

	err := query.First(&testing, "id = ?", form.ID).Error
	if err != nil {
		error2.ErrXtremeTestingGet(err.Error())
	}

	return testing
}

func (repo *testingRepository) FindByForm(form form.TestingFilterForm) []model.Testing {
	query := repo.prepareAndQuery(form)

	var testings []model.Testing
	err := query.Find(&testings).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		error2.ErrXtremeTestingGet(err.Error())
	}

	return testings
}

func (repo *testingRepository) PaginateByForm(form form.TestingFilterForm) ([]model.Testing, interface{}) {
	if form.Page <= 0 {
		form.Page = 1
	}

	parameter := url.Values{}
	parameter.Set("page", strconv.Itoa(form.Page))
	parameter.Set("limit", strconv.Itoa(form.Limit))

	query := repo.prepareAndQuery(form)

	testings, pagination, err := xtrememodel.Paginate(query, parameter, model.Testing{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		error2.ErrXtremeTestingGet(err.Error())
	}

	return testings, pagination
}

func (repo *testingRepository) Create(form form.TestingForm) model.Testing {
	testing := model.Testing{
		Name: form.Name,
	}

	err := repo.transaction.Create(&testing).Error
	if err != nil {
		error2.ErrXtremeTestingSave(err.Error())
	}

	return testing
}

func (repo *testingRepository) Delete(testing model.Testing) {
	err := repo.transaction.Delete(&testing).Error
	if err != nil {
		error2.ErrXtremeTestingDelete(err.Error())
	}
}

func (repo *testingRepository) AddSub(testing model.Testing, sub string) model.TestingSub {
	testingSub := model.TestingSub{
		TestingId: testing.ID,
		Name:      sub,
	}

	err := repo.transaction.Create(&testingSub).Error
	if err != nil {
		error2.ErrXtremeTestingSubSave(err.Error())
	}

	return testingSub
}

func (repo *testingRepository) DeleteSub(testingSub model.TestingSub) {
	err := repo.transaction.Delete(&testingSub).Error
	if err != nil {
		error2.ErrXtremeTestingSubDelete(err.Error())
	}
}

/** --- UNEXPORTED FUNCTIONS --- */

func (repo *testingRepository) prepareAndQuery(form form.TestingFilterForm) *gorm.DB {
	var query *gorm.DB
	if form.UseTransaction {
		query = repo.transaction
	} else {
		query = config.PgSQL
	}

	if form.FromDate != "" && form.ToDate != "" {
		fromDate, toDate := core.SetDateRange(form.FromDate, form.ToDate)

		query = query.Preload("Subs").
			Where(`"createdAt" BETWEEN ? AND ?`, fromDate, toDate)
	}

	if search := form.Search; len(search) > 3 {
		searchVal := "%" + search + "%"
		query = query.Where("name ILIKE ?", searchVal)
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
