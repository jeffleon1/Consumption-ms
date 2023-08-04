package infraestructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeffleon1/consumption-ms/pkg/application"
	"github.com/jeffleon1/consumption-ms/pkg/application/applicationfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

const (
	ErrorWhenDataIsWrong       = "when data is wrong in request"
	ErrorShouldReturnaE        = "should return an error"
	ConsumptionPath            = "/consumption"
	ConsumptionInformationPath = "/consumption/information"
)

var _ = Describe("CreatePowerConsumptionRecords", func() {
	var (
		router                      *gin.Engine
		server                      *ghttp.Server
		mockPowerConsumptionService *applicationfakes.FakePowerConsumptionService
	)

	BeforeEach(func() {
		router = gin.Default()
		mockPowerConsumptionService = &applicationfakes.FakePowerConsumptionService{}
		mockHandler := NewPowerConsumptionHandler(mockPowerConsumptionService)
		router.GET(ConsumptionPath, mockHandler.GetConsumptionByMeterIDAndWindowTime)
		server = ghttp.NewServer()
		server.RouteToHandler("GET", ConsumptionPath, router.ServeHTTP)
	})

	AfterEach(func() {
		server.Close()
	})

	Context(ErrorWhenDataIsWrong, func() {
		It(ErrorShouldReturnaE, func() {
			mockPowerConsumptionService.GetConsumptionByMeterIDAndWindowTimeReturns(nil, fmt.Errorf("Some error"))
			resp, err := http.Get(fmt.Sprintf("%s%s", server.URL(), fmt.Sprintf("%s?meter_ids=1,2&start_date=2023-05-30&end_date=2023-07-01&kind_period=weekly", ConsumptionPath)))
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

			var responseBody Response
			json.NewDecoder(resp.Body).Decode(&responseBody)
			Expect(responseBody.Status).To(Equal("ERROR"))
			Expect(responseBody.Msg).To(ContainSubstring("Something goes wrong"))
			Expect(mockPowerConsumptionService.GetConsumptionByMeterIDAndWindowTimeCallCount()).To(Equal(1))
		})
	})

	Context(ErrorWhenDataIsWrong, func() {
		It(ErrorShouldReturnaE, func() {
			mockPowerConsumptionService.GetConsumptionByMeterIDAndWindowTimeReturns(nil, fmt.Errorf("Some error"))
			resp, err := http.Get(fmt.Sprintf("%s%s", server.URL(), ConsumptionPath))
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

			var responseBody Response
			json.NewDecoder(resp.Body).Decode(&responseBody)
			Expect(responseBody.Status).To(Equal("ERROR"))
			Expect(responseBody.Msg).To(ContainSubstring("Something goes wrong"))

			Expect(mockPowerConsumptionService.GetConsumptionByMeterIDAndWindowTimeCallCount()).To(Equal(0))
		})
	})

	Context("when data is OK in request", func() {
		It("should return success response", func() {
			serializers := []application.Serializer{
				{
					Period:             []string{"Jun 19", "Jul 19", "Ago 19"},
					MeterID:            1,
					Active:             []float64{100.0, 120.5, 80.2},
					ReactiveInductive:  []float64{50.2, 40.5, 60.0},
					ReactiveCapacitive: []float64{30.0, 35.2, 28.7},
					Exported:           []float64{20.5, 25.8, 18.6},
				},
				{
					Period:             []string{"Sep 19", "Oct 19", "Nov 19"},
					MeterID:            2,
					Active:             []float64{90.2, 110.7, 95.3},
					ReactiveInductive:  []float64{48.6, 42.0, 55.1},
					ReactiveCapacitive: []float64{29.8, 31.5, 26.4},
					Exported:           []float64{18.7, 22.3, 20.0},
				},
			}

			mockPowerConsumptionService.GetConsumptionByMeterIDAndWindowTimeReturns(serializers, nil)
			resp, err := http.Get(fmt.Sprintf("%s%s", server.URL(), fmt.Sprintf("%s?meter_ids=1,2&start_date=2023-05-30&end_date=2023-07-01&kind_period=weekly", ConsumptionPath)))
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var responseBody Response
			json.NewDecoder(resp.Body).Decode(&responseBody)
			Expect(responseBody.Status).To(Equal("SUCCESS"))
			Expect(responseBody.Msg).To(ContainSubstring("information successfully brought"))

			Expect(mockPowerConsumptionService.GetConsumptionByMeterIDAndWindowTimeCallCount()).To(Equal(1))
		})
	})
})

var _ = Describe("ImportCsvToDatabase", func() {
	var (
		router                      *gin.Engine
		server                      *ghttp.Server
		mockPowerConsumptionService *applicationfakes.FakePowerConsumptionService
	)

	BeforeEach(func() {
		router = gin.Default()
		mockPowerConsumptionService = &applicationfakes.FakePowerConsumptionService{}
		mockHandler := NewPowerConsumptionHandler(mockPowerConsumptionService)
		router.POST(ConsumptionInformationPath, mockHandler.GetConsumptionByMeterIDAndWindowTime)
		server = ghttp.NewServer()
		server.RouteToHandler("POST", ConsumptionInformationPath, router.ServeHTTP)
	})

	AfterEach(func() {
		server.Close()
	})

	Context(ErrorWhenDataIsWrong, func() {
		It(ErrorShouldReturnaE, func() {
			var buf bytes.Buffer
			multipartWriter := multipart.NewWriter(&buf)
			fileWriter, err := multipartWriter.CreateFormFile("file", "example.csv")
			mockPowerConsumptionService.ImportCsvToDatabaseReturns(nil)
			Expect(err).To(BeNil())
			_, err = io.Copy(fileWriter, bytes.NewBufferString("ID,MeterID,ActiveEnergy,ReactiveEnergy,CapacitiveReactive,Solar,Date\n1,2,100,50,30,20,2023-08-01 12:00:00\n"))
			Expect(err).To(BeNil())
			multipartWriter.Close()

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", ConsumptionInformationPath),
				ghttp.VerifyContentType("multipart/form-data"),
				ghttp.RespondWith(http.StatusCreated, nil),
			))

			req, err := http.NewRequest("POST", server.URL()+ConsumptionInformationPath, &buf)
			Expect(err).To(BeNil())
			req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

			resp, err := http.DefaultClient.Do(req)
			var responseBody Response
			json.NewDecoder(resp.Body).Decode(&responseBody)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})
	})

})
