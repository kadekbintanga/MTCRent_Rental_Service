package model

import "gorm.io/gorm"

type Customer struct {
	ID        uint           `gorm:"column:id"`
	UUID      string         `gorm:"column:uuid;type:char(36)"`
	Name      string         `gorm:"column:name;type:varchar(250);not null"`
	IDNumber  string         `gorm:"column:IDNumber;type:varchar(100);not null"`
	SIMNumber string         `gorm:"column:SIMNumber;type:varchar(100);not null"`
	Phone     string         `gorm:"column:phone;type:varchar(30);not null"`
	StatusId  int            `gorm:"column:statusId"`
	DeletedAt gorm.DeletedAt `gorm:"column:deletedAt;index"`
}

func (Customer) TableName() string {
	return "customers"
}

func (model Customer) SetReference() uint {
	return model.ID
}
