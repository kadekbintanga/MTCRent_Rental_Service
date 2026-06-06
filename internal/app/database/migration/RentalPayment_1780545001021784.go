package migration

import (
	"os"
	"service/internal/pkg/config"
	"service/internal/pkg/model"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
)

type RentalPayment_1780545001021784 struct{}

func (RentalPayment_1780545001021784) Reference() string {
	return "RentalPayment_1780545001021784"
}

func (RentalPayment_1780545001021784) Tables() []xtremedb.Table {
	owner := os.Getenv("DB_OWNER")
	return []xtremedb.Table{
		{Connection: config.PgSQL, CreateTable: model.RentalPayment{}, Owner: owner},
	}
}

func (RentalPayment_1780545001021784) Columns() []xtremedb.Column {
	return []xtremedb.Column{}
}
