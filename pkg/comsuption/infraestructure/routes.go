package infraestructure

import "github.com/gin-gonic/gin"

type PowerConsumptionRoutes struct {
	powerConsumptionHandler *PowerConsumptionHandlerImpl
}

func (ro *PowerConsumptionRoutes) RegisterRoutes(public *gin.RouterGroup) {
	public.GET("/example", ro.powerConsumptionHandler.GetConsumptionByMeterIDAndWindowTime)
	public.POST("/", ro.powerConsumptionHandler.ImportCsvToDatabase)
}

func NewRoutes(powerConsumptionHandler *PowerConsumptionHandlerImpl) *PowerConsumptionRoutes {
	return &PowerConsumptionRoutes{
		powerConsumptionHandler,
	}
}
