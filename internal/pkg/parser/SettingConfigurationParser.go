package parser

import "service/internal/pkg/model"

type SettingConfigurationParser struct {
	Array  []model.SettingConfiguration
	Object model.SettingConfiguration
}

func (parser SettingConfigurationParser) Get() []interface{} {
	var result []interface{}

	for _, settingConfig := range parser.Array {
		firstPaser := SettingConfigurationParser{Object: settingConfig}
		result = append(result, firstPaser.First())
	}
	return result
}

func (parser SettingConfigurationParser) First() interface{} {
	settingConfig := parser.Object

	return map[string]interface{}{
		"key":       settingConfig.Key,
		"value":     settingConfig.Value,
		"createdAt": settingConfig.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt": settingConfig.UpdatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser SettingConfigurationParser) CreateActivity(action string) interface{} {
	settingConfig := parser.Object

	return map[string]interface{}{
		"id":        settingConfig.ID,
		"name":      settingConfig.Key,
		"createdAt": settingConfig.CreatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser SettingConfigurationParser) UpdateActivity(action string) interface{} {
	return parser.CreateActivity(action)
}

func (parser SettingConfigurationParser) DeleteActivity(action string) interface{} {
	return parser.CreateActivity(action)
}

func (parser SettingConfigurationParser) GeneralActivity(action string) interface{} {
	if action == "onlyName" {
		settingConfig := parser.Object

		return map[string]interface{}{
			"name": settingConfig.Key,
		}
	}

	return parser.CreateActivity(action)
}
