package migration

import (
	"os"
	"service/internal/pkg/config"
	"service/internal/pkg/model"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
)

type Customer_1780475661082866 struct{}

func (Customer_1780475661082866) Reference() string {
	return "Customer_1780475661082866"
}

func (Customer_1780475661082866) Tables() []xtremedb.Table {
	owner := os.Getenv("DB_OWNER")
	return []xtremedb.Table{
		{Connection: config.PgSQL, CreateTable: model.Customer{}, Owner: owner},
	}
}

func (Customer_1780475661082866) Columns() []xtremedb.Column {
	return []xtremedb.Column{}
}
