package seeder

import (
	"service/internal/pkg/config"
	"service/internal/pkg/model"
)

type MotorcycleComponentBrandDefault struct{}

func (seed *MotorcycleComponentBrandDefault) Seed() {
	motorcycleBrands := seed.setMotorcycleComponentBrandDefaultData()
	for _, motorcycleBrand := range motorcycleBrands {
		var count int64
		config.PgSQL.Model(&model.MotorcycleComponentBrand{}).Where(`motorcycle_component_brands."name" = ?`, motorcycleBrand["name"]).Count(&count)
		if count > 0 {
			continue
		}

		config.PgSQL.Create((&model.MotorcycleComponentBrand{
			Name:    motorcycleBrand["name"].(string),
			Default: motorcycleBrand["default"].(bool),
		}))
	}
}

/** --- UNEXPORTED FUNCTIONS --- */

func (seed *MotorcycleComponentBrandDefault) setMotorcycleComponentBrandDefaultData() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":    "OTHER",
			"default": true,
		},
	}
}
