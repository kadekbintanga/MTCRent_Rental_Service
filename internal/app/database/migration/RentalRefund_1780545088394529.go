package migration

import (
	"os"
	"service/internal/pkg/config"
	"service/internal/pkg/model"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
)

type RentalRefund_1780545088394529 struct{}

func (RentalRefund_1780545088394529) Reference() string {
	return "RentalRefund_1780545088394529"
}

func (RentalRefund_1780545088394529) Tables() []xtremedb.Table {
	owner := os.Getenv("DB_OWNER")
	return []xtremedb.Table{
		{Connection: config.PgSQL, CreateTable: model.RentalRefund{}, Owner: owner},
	}
}

func (RentalRefund_1780545088394529) Columns() []xtremedb.Column {
	return []xtremedb.Column{}
}
