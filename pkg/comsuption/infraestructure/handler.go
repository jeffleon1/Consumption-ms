package infraestructure

import "github.com/jeffleon1/consumption-ms/pkg/comsuption/application"

type PowerConsumptionHandlerImpl interface {
}

type PowerConsumptionHandler struct {
	powerConsumptionService *application.PowerConsumptionService
}

func NewPowerConsumptionHandler(powerConsumptionService *application.PowerConsumptionService) PowerConsumptionHandlerImpl {
	return &PowerConsumptionHandler{
		powerConsumptionService,
	}
}
