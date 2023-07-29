package repositories

import "github.com/jeffleon1/consumption-ms/pkg/comsuption/domain"

type PostgreSQLPowerConsumptionRepository struct {
}

func NewPostgreSQLPowerConsumptionRepository() domain.PowerConsumptionRepository {
	return &PostgreSQLPowerConsumptionRepository{}
}
