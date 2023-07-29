package domain

import "time"

type UserConsumption struct {
	ID                 string    `json:"id"`
	MeterID            string    `json:"meter_id"`
	ActiveEnergy       float64   `json:"active_energy"`
	ReactiveEnergy     float64   `json:"reactive_energy"`
	CapacitiveReactive float64   `json:"capacity_energy"`
	Solar              float64   `json:"solar"`
	Date               time.Time `json:"date"`
}

type PowerConsumptionRepository interface {
}
