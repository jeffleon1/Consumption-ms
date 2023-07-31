package infraestructure

import (
	"fmt"
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
	meterIDs := c.Query("meters_ids")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	kindPeriod := c.Query("kind_period")
	if meterIDs == "" || startDate == "" || endDate == "" || kindPeriod == "" {
		c.JSON(http.StatusBadRequest, Response{
			Msg:    "Something goes wrong with your query params",
			Status: "ERROR",
			Data:   nil,
			Err:    fmt.Sprintf("Some params are blank meter_ids=%s start_date=%s end_date=%s kind_period=%s", meterIDs, startDate, endDate, kindPeriod),
		})
		return
	}

	s.powerConsumptionService.GetConsumptionByMeterIDAndWindowTime(meterIDs, startDate, endDate, kindPeriod)

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
	c.JSON(http.StatusCreated, Response{
		Msg:    "All records were successfully saved",
		Status: "SUCCESS",
		Data:   nil,
		Err:    nil,
	})
}
