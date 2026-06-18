package parser

import "service/internal/pkg/model"

type MotorcycleBrandParser struct {
	Array  []model.MotorcycleComponentBrand
	Object model.MotorcycleComponentBrand
}

func (parser MotorcycleBrandParser) Get() []interface{} {
	var result []interface{}

	for _, motorcycleBrand := range parser.Array {
		firstParser := MotorcycleBrandParser{Object: motorcycleBrand}
		result = append(result, firstParser.First())
	}
	return result
}

func (parser MotorcycleBrandParser) First() interface{} {
	motorcycleBrand := parser.Object

	return map[string]interface{}{
		"id":        motorcycleBrand.ID,
		"name":      motorcycleBrand.Name,
		"createdAt": motorcycleBrand.CreatedAt.Format("02/01/2006 15:04"),
		"updatedAt": motorcycleBrand.UpdatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser MotorcycleBrandParser) CreateActivity(action string) interface{} {
	motorcycleBrand := parser.Object

	return map[string]interface{}{
		"id":        motorcycleBrand.ID,
		"name":      motorcycleBrand.Name,
		"createdAt": motorcycleBrand.CreatedAt.Format("02/01/2006 15:04"),
	}
}

func (parser MotorcycleBrandParser) UpdateActivity(action string) interface{} {
	return parser.CreateActivity(action)
}

func (parser MotorcycleBrandParser) DeleteActivity(action string) interface{} {
	return parser.CreateActivity(action)
}

func (parser MotorcycleBrandParser) GeneralActivity(action string) interface{} {
	return parser.CreateActivity(action)
}
