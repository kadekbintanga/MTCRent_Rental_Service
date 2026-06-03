package database

import (
	xtremedb "github.com/globalxtreme/go-core/v2/database"

	"service/internal/app/database/migration"
)

func Migrations() []xtremedb.Migration {
	return []xtremedb.Migration{
		&migration.Activity_1726651211960757{},
		&migration.MotorcycleBrand_1779349665811637{},
		&migration.Motorcycle_1779349674123225{},
	}
}
