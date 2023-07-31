package infraestructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeffleon1/consumption-ms/pkg/comsuption/application"
)

type Response struct {
	Msg    string      `json:"msg"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Err    interface{} `json:"error"`
}

type PowerConsumptionHandlerImpl struct {
	powerConsumptionService application.PowerConsumptionService
}

func NewPowerConsumptionHandler(powerConsumptionService application.PowerConsumptionService) *PowerConsumptionHandlerImpl {
	return &PowerConsumptionHandlerImpl{
		powerConsumptionService,
	}
}

func (s *PowerConsumptionHandlerImpl) GetConsumptionByMeterIDAndWindowTime(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (s *PowerConsumptionHandlerImpl) ImportCsvToDatabase(c *gin.Context) {
	csvPartFile, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Msg:    "Something goes wrong please check your csv file",
			Status: "ERROR",
			Data:   nil,
			Err:    err.Error(),
		})
		return
	}
	err = s.powerConsumptionService.ImportCsvToDatabase(&csvPartFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Msg:    "Something goes wrong please check your csv file",
			Status: "ERROR",
			Data:   nil,
			Err:    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
