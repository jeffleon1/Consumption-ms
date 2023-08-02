package infraestructure

import "github.com/gin-gonic/gin"

type PowerConsumptionRoutes struct {
	powerConsumptionHandler *PowerConsumptionHandlerImpl
}

func (ro *PowerConsumptionRoutes) RegisterRoutes(public *gin.RouterGroup) {
	public.GET("/consumption", ro.powerConsumptionHandler.GetConsumptionByMeterIDAndWindowTime)
	public.POST("/consumption/information", ro.powerConsumptionHandler.ImportCsvToDatabase)
}

func NewRoutes(powerConsumptionHandler *PowerConsumptionHandlerImpl) *PowerConsumptionRoutes {
	return &PowerConsumptionRoutes{
		powerConsumptionHandler,
	}
}
