package migration

import (
	"os"

	xtremedb "github.com/globalxtreme/go-core/v2/database"

	"service/internal/pkg/config"
	"service/internal/pkg/model"
)

type SettingConfiguration_1780466563416432 struct{}

func (SettingConfiguration_1780466563416432) Reference() string {
	return "SettingConfiguration_1780466563416432"
}

func (SettingConfiguration_1780466563416432) Tables() []xtremedb.Table {
	owner := os.Getenv("DB_OWNER")
	return []xtremedb.Table{
		{Connection: config.PgSQL, CreateTable: model.SettingConfiguration{}, Owner: owner},
	}
}

func (SettingConfiguration_1780466563416432) Columns() []xtremedb.Column {
	return []xtremedb.Column{}
}
