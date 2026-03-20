package service

import (
	"fmt"
	xtremefs "github.com/globalxtreme/go-core/v2/filesystem"
	"gorm.io/gorm"
	"service/internal/pkg/activity"
	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	error2 "service/internal/pkg/error"
	form2 "service/internal/pkg/form"
	"service/internal/pkg/model"
	"service/internal/pkg/port"
	"service/internal/testing/repository"
)

// TODO: Hanya contoh. nanti langsung hapus saja
type TestingService interface {
	SetTransaction(tx *gorm.DB)
	SetActivityRepository(repo port.ActivityRepository)

	Create(form form2.TestingForm) model.Testing
	CreateConsumer(form form2.TestingForm) model.Testing
	UploadByFile(form form2.TestingUploadForm) map[string]interface{}
	UploadByContent(form form2.TestingUploadContentForm) map[string]interface{}
}

func NewTestingService() TestingService {
	return &testingService{}
}

type testingService struct {
	tx *gorm.DB

	repository   repository.TestingRepository
	activityRepo port.ActivityRepository
}

func (srv *testingService) SetTransaction(tx *gorm.DB) {
	srv.tx = tx
}

func (srv *testingService) SetActivityRepository(repo port.ActivityRepository) {
	srv.activityRepo = repo
}

func (srv *testingService) Create(form form2.TestingForm) model.Testing {
	var testing model.Testing

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewTestingRepository(tx)

		testing = srv.repository.Create(form)

		for _, sub := range form.Subs {
			testingSub := srv.repository.AddSub(testing, sub)
			testing.Subs = append(testing.Subs, testingSub)
		}

		activity.UseActivity{}.SetReference(testing).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Enter new testing: %s [%d]", testing.Name, testing.ID))

		return nil
	})

	return testing
}

func (srv *testingService) CreateConsumer(form form2.TestingForm) model.Testing {
	var testing model.Testing

	config.PgSQL.Transaction(func(tx *gorm.DB) error {
		srv.repository = repository.NewTestingRepository(tx)

		testing = srv.repository.Create(form)

		for _, sub := range form.Subs {
			testingSub := srv.repository.AddSub(testing, sub)
			testing.Subs = append(testing.Subs, testingSub)
		}

		activity.UseActivity{}.SetReference(testing).SetNewProperty(constant.ACTION_CREATE).
			Save(fmt.Sprintf("Enter new testing: %s [%d]", testing.Name, testing.ID))

		return nil
	})

	return testing
}

func (srv *testingService) UploadByFile(form form2.TestingUploadForm) map[string]interface{} {
	uploader := xtremefs.Uploader{Path: constant.PathImageTesting(), IsPublic: true}
	filePath, err := uploader.MoveFile(form.Request, "testFile[testing][0]")
	if err != nil {
		error2.ErrXtremeTestingSave("Unable to upload file: " + err.Error())
	}

	storage := xtremefs.Storage{IsPublic: uploader.IsPublic}

	return map[string]interface{}{
		"url":      storage.GetFullPathURL(filePath.(string)),
		"fullPath": storage.GetFullPath(filePath.(string)),
		"path":     filePath.(string),
	}
}

func (srv *testingService) UploadByContent(form form2.TestingUploadContentForm) map[string]interface{} {
	uploader := xtremefs.Uploader{Path: constant.PathImageTesting(), IsPublic: true}
	filePath, err := uploader.MoveContent(form.Content)
	if err != nil {
		error2.ErrXtremeTestingSave("Unable to upload file: " + err.Error())
	}

	storage := xtremefs.Storage{IsPublic: uploader.IsPublic}

	return map[string]interface{}{
		"url":      storage.GetFullPathURL(filePath.(string)),
		"fullPath": storage.GetFullPath(filePath.(string)),
		"path":     filePath.(string),
	}
}
