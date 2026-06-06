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
		&migration.SettingConfiguration_1780466563416432{},
		&migration.Customer_1780475661082866{},
		&migration.Rental_1780545147162427{},
		&migration.RentalPayment_1780545001021784{},
		&migration.RentalRefund_1780545088394529{},
	}
}
