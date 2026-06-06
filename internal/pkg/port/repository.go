package port

import (
	"net/url"

	"service/internal/pkg/core"
	"service/internal/pkg/form"
	"service/internal/pkg/form/option"
	"service/internal/pkg/model"
)

/** --- ACTIVITY --- */

type ActivityRepository interface {
	Find(parameters url.Values) ([]model.Activity, interface{}, error)
}

type CustomerRepository interface {
	core.TransactionInterface
	core.FirstRepository[option.CustomerOption, model.Customer]

	Create(opt option.CustomerSaveOption) model.Customer
}

type MotorcycleRepository interface {
	core.TransactionInterface
	core.FindRepository[form.MotorcycleFilterForm, model.Motorcycle]

	UpdateStatus(motorcycle model.Motorcycle, form form.MotorcycleStatusUpdateForm) model.Motorcycle
}
