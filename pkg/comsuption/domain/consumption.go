package domain

import (
	"mime/multipart"

	"gorm.io/gorm"
)

type UserConsumption struct {
	gorm.Model
	ID                 string  `gorm:"primary_key;auto_increment" json:"id" csv:"id"`
	MeterID            string  `gorm:"meter_id" json:"meter_id" csv:"meter_id"`
	ActiveEnergy       float64 `gorm:"active_energy" json:"active_energy" csv:"active_energy"`
	ReactiveEnergy     float64 `gorm:"reactive_energy" json:"reactive_energy" csv:"reactive_energy"`
	CapacitiveReactive float64 `gorm:"capacity_energy" json:"capacitive_reactive" csv:"capacitive_reactive"`
	Solar              float64 `gorm:"solar" json:"solar" csv:"solar"`
	Date               string  `gorm:"date" json:"date" csv:"date"`
}

type PostgresPowerConsumptionRepository interface {
	GetConsumptionByMeterIDAndWindowTime(startDate string, endDate string) ([]UserConsumption, error)
	CreatePowerConsumptionRecords(usersPowerConsumption []*UserConsumption) error
	ModelMigration() error
}

type CSVPowerConsumptionRepository interface {
	ConvertCSVToStruct(file *multipart.File) ([]*UserConsumption, error)
}
