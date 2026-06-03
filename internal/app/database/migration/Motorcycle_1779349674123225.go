package migration

import (
	"os"

	xtremedb "github.com/globalxtreme/go-core/v2/database"

	"service/internal/pkg/config"
	"service/internal/pkg/model"
)

type Motorcycle_1779349674123225 struct{}

func (Motorcycle_1779349674123225) Reference() string {
	return "Motorcycle_1779349674123225"
}

func (Motorcycle_1779349674123225) Tables() []xtremedb.Table {
	owner := os.Getenv("DB_OWNER")
	return []xtremedb.Table{
		{Connection: config.PgSQL, CreateTable: model.Motorcycle{}, Owner: owner},
	}
}

func (Motorcycle_1779349674123225) Columns() []xtremedb.Column {
	return []xtremedb.Column{}
}
