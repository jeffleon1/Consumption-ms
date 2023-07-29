package infraestructure

import "github.com/gin-gonic/gin"

type PowerConsumptionRoutes struct {
	powerConsumptionHandler *PowerConsumptionHandlerImpl
}

func (ro *PowerConsumptionRoutes) RegisterRoutes(public *gin.RouterGroup) {
	//Implementation soon
}

func NewRoutes(powerConsumptionHandler *PowerConsumptionHandlerImpl) *PowerConsumptionRoutes {
	return &PowerConsumptionRoutes{
		powerConsumptionHandler,
	}
}
