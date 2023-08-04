package infraestructure

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeffleon1/consumption-ms/pkg/application"
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

// Get the user consumption information in a window time divided monthly, weekly or daily
// @Tags Consumption
// @Summary Get the user consumption information in a window time divided monthly, weekly or daily
// @Description Get the user consumption information in a window time divided monthly, weekly or daily
// @Accept  json
// @Produce  json
// @Param start_date query string  true  "start date"
// @Param end_date query string  true  "end date"
// @Param kind_period query string  true  "kind period"
// @Param meter_ids query string  true "meter ids"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /consumption [get]
func (s *PowerConsumptionHandlerImpl) GetConsumptionByMeterIDAndWindowTime(c *gin.Context) {
	meterIDs := c.Query("meter_ids")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	kindPeriod := c.Query("kind_period")
	if meterIDs == "" || startDate == "" || endDate == "" || kindPeriod == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, Response{
			Msg:    "Something goes wrong with your query params",
			Status: "ERROR",
			Data:   nil,
			Err:    fmt.Sprintf("Some params are blank meter_ids=%s start_date=%s end_date=%s kind_period=%s", meterIDs, startDate, endDate, kindPeriod),
		})
		return
	}

	filterSerializer := &FilterConsumptionSerializer{}

	data, err := s.powerConsumptionService.GetConsumptionByMeterIDAndWindowTime(meterIDs, startDate, endDate, kindPeriod)
	if err != nil {
		fmt.Println("Entro aca con todos los poderes")
		c.AbortWithStatusJSON(http.StatusBadRequest, Response{
			Msg:    "Something goes wrong",
			Status: "ERROR",
			Data:   nil,
			Err:    err.Error(),
		})
		return
	}

	filterSerializer.ToFilterConsumptionSerializer(data)

	c.JSON(http.StatusOK, Response{
		Msg:    "information successfully brought",
		Status: "SUCCESS",
		Data:   filterSerializer,
		Err:    nil,
	})
}

// Import a csv file to insert the information in the user_consumption database
// @Tags Consumption
// @Summary Import a csv file to insert the information in the user_consumption database
// @Description Import a csv file to insert the information in the user_consumption database
// @Accept  json
// @Produce  json
// @Param file	formData file true "this is a csv test file"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /consumption/information [post]
func (s *PowerConsumptionHandlerImpl) ImportCsvToDatabase(c *gin.Context) {
	csvPartFile, _, err := c.Request.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, Response{
			Msg:    "Something goes wrong please check your csv file",
			Status: "ERROR",
			Data:   nil,
			Err:    err.Error(),
		})
		return
	}
	err = s.powerConsumptionService.ImportCsvToDatabase(&csvPartFile)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, Response{
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
