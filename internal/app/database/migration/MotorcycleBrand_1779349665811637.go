package migration

import (
	"os"

	xtremedb "github.com/globalxtreme/go-core/v2/database"

	"service/internal/pkg/config"
	"service/internal/pkg/model"
)

type MotorcycleBrand_1779349665811637 struct{}

func (MotorcycleBrand_1779349665811637) Reference() string {
	return "MotorcycleBrand_1779349665811637"
}

func (MotorcycleBrand_1779349665811637) Tables() []xtremedb.Table {
	owner := os.Getenv("DB_OWNER")
	return []xtremedb.Table{
		{Connection: config.PgSQL, CreateTable: model.MotorcycleComponentBrand{}, Owner: owner},
	}
}

func (MotorcycleBrand_1779349665811637) Columns() []xtremedb.Column {
	return []xtremedb.Column{}
}
