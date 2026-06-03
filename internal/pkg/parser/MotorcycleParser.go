package parser

import (
	"service/internal/pkg/constant"
	"service/internal/pkg/model"
)

type MotorcycleParser struct {
	Array  []model.Motorcycle
	Object model.Motorcycle
}

func (parser MotorcycleParser) Get() []interface{} {
	var result []interface{}

	for _, motorcycle := range parser.Array {
		firstParser := MotorcycleParser{Object: motorcycle}
		result = append(result, firstParser.First())
	}
	return result
}

func (parser MotorcycleParser) First() interface{} {
	motorcycle := parser.Object

	return map[string]interface{}{
		"uuid":        motorcycle.UUID,
		"plateNumber": motorcycle.PlateNumber,
		"name":        motorcycle.Name,
		"typeId":      constant.MotorcycleType{}.IDAndName(motorcycle.TypeId),
		"year":        motorcycle.Year,
		"pricePerDay": motorcycle.PricePerDay,
		"statusId":    constant.MotorcycleStatus{}.IDAndName(motorcycle.StatusId),
		"brandId":     motorcycle.BrandId,
		"brandName":   motorcycle.Brand.Name,
		"createdAt":   motorcycle.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt":   motorcycle.UpdatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser MotorcycleParser) CreateActivity(action string) interface{} {
	motocycle := parser.Object

	return map[string]interface{}{
		"id":        motocycle.ID,
		"name":      motocycle.Name,
		"createdAt": motocycle.CreatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser MotorcycleParser) UpdateActivity(action string) interface{} {
	return parser.CreateActivity(action)
}

func (parser MotorcycleParser) DeleteActivity(action string) interface{} {
	return parser.CreateActivity(action)
}

func (parser MotorcycleParser) GeneralActivity(action string) interface{} {
	if action == "onlyName" {
		motorcycleBrand := parser.Object

		return map[string]interface{}{
			"name": motorcycleBrand.Name,
		}
	}

	return parser.CreateActivity(action)
}
