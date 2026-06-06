package seeder

import (
	"service/internal/pkg/config"
	"service/internal/pkg/model"
)

type SettingConfiguration struct{}

func (seed *SettingConfiguration) Seed() {
	settingConfigs := seed.setSettingConfigurationData()
	for _, settingConfig := range settingConfigs {
		var count int64
		config.PgSQL.Model(&model.SettingConfiguration{}).Where(`setting_configurations."key" = ?`, settingConfig["key"]).Count(&count)
		if count > 0 {
			continue
		}

		config.PgSQL.Create(&model.SettingConfiguration{
			Key:   settingConfig["key"].(string),
			Value: settingConfig["value"].(string),
		})
	}
}

/** --- UNEXPORTED FUNCTIONS --- */

func (seed *SettingConfiguration) setSettingConfigurationData() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"key":   "blacklistLimitDay",
			"value": "5",
		},
	}
}
