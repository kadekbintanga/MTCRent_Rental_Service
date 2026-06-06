package migration

import (
	"os"
	"service/internal/pkg/config"
	"service/internal/pkg/model"

	xtremedb "github.com/globalxtreme/go-core/v2/database"
)

type Rental_1780545147162427 struct{}

func (Rental_1780545147162427) Reference() string {
	return "Rental_1780545147162427"
}

func (Rental_1780545147162427) Tables() []xtremedb.Table {
	owner := os.Getenv("DB_OWNER")
	return []xtremedb.Table{
		{Connection: config.PgSQL, CreateTable: model.Rental{}, Owner: owner},
	}
}

func (Rental_1780545147162427) Columns() []xtremedb.Column {
	return []xtremedb.Column{}
}
