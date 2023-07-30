package application

import "github.com/jeffleon1/consumption-ms/pkg/comsuption/domain"

type PowerConsumptionService interface {
}

type PowerConsumptionServiceImpl struct {
	postgresRepository *domain.PowerConsumptionRepository
}

func NewPowerConsumptionService(postgresRepository *domain.PowerConsumptionRepository) PowerConsumptionService {
	return &PowerConsumptionServiceImpl{
		postgresRepository,
	}
}
