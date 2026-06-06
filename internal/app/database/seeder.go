package database

import "service/internal/app/database/seeder"

type DatabaseSeeder interface {
	Seed()
}

type data struct {
	DatabaseSeeder
}

func Seeder() {
	seeders := []data{
		{&seeder.TestingSeeder{}},
		{&seeder.SettingConfiguration{}},
		{&seeder.MotorcycleComponentBrandDefault{}},
	}

	for _, seed := range seeders {
		seed.Seed()
	}
}
