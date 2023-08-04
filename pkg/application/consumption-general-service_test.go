package application

import (
	"errors"
	"fmt"
	"time"

	"github.com/jeffleon1/consumption-ms/pkg/domain"
	"github.com/jeffleon1/consumption-ms/pkg/domain/domainfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	startDate string = "2023-01-01"
	endDate   string = "2023-02-01"
)

var _ = Describe("chekingQueryParamConstrains", func() {

	var (
		mockMySQLRepo               *domainfakes.FakeMySQLPowerConsumptionRepository
		mockCSVRepo                 *domainfakes.FakeCSVPowerConsumptionRepository
		mockPowerConsumptionService PowerConsumptionService
	)

	BeforeEach(func() {
		mockMySQLRepo = &domainfakes.FakeMySQLPowerConsumptionRepository{}
		mockCSVRepo = &domainfakes.FakeCSVPowerConsumptionRepository{}
		mockPowerConsumptionService = NewPowerConsumptionService(mockMySQLRepo, mockCSVRepo)
	})

	Context("checkingQueryParamConstrains", func() {
		It("should return correct UserConsumptionQueryParams for valid parameters", func() {
			meterIDs := "1,2,3"
			kindPeriod := "monthly"

			queryParams, err := mockPowerConsumptionService.CheckingQueryParamConstrains(meterIDs, kindPeriod, startDate, endDate)

			Expect(err).To(BeNil())
			Expect(queryParams).ToNot(BeNil())
			Expect(queryParams.StartDate).To(Equal(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))
			Expect(queryParams.EndDate).To(Equal(time.Date(2023, 2, 1, 23, 59, 59, 0, time.UTC)))
			Expect(queryParams.MeterIDs).To(Equal([]int{1, 2, 3}))
			Expect(queryParams.KindPeriod).To(Equal("monthly"))
		})

		It("should handle invalid startDate format", func() {
			meterIDs := "1,2,3"
			kindPeriod := "monthly"
			invalidStartDate := "2023/04-02"
			queryParams, err := mockPowerConsumptionService.CheckingQueryParamConstrains(meterIDs, kindPeriod, invalidStartDate, endDate)

			Expect(err).To(HaveOccurred())
			Expect(queryParams).To(BeNil())
		})

		It("should handle empty meterIDs correctly", func() {
			meterIDs := ""
			kindPeriod := "monthly"

			queryParams, err := mockPowerConsumptionService.CheckingQueryParamConstrains(meterIDs, kindPeriod, startDate, endDate)
			Expect(err).To(HaveOccurred())
			Expect(queryParams).To(BeNil())
		})

		It("should handle invalid kindPeriod correctly", func() {
			meterIDs := "1,2,3"
			kindPeriod := "invalid"
			queryParams, err := mockPowerConsumptionService.CheckingQueryParamConstrains(meterIDs, kindPeriod, startDate, endDate)
			fmt.Println(queryParams)
			fmt.Println(err)
			Expect(err).To(HaveOccurred())
			Expect(queryParams).To(BeNil())
		})
	})

})

var _ = Describe("PowerConsumptionService", func() {
	var (
		mockMySQLRepo               *domainfakes.FakeMySQLPowerConsumptionRepository
		mockCSVRepo                 *domainfakes.FakeCSVPowerConsumptionRepository
		mockPowerConsumptionService PowerConsumptionService
	)

	BeforeEach(func() {
		mockMySQLRepo = &domainfakes.FakeMySQLPowerConsumptionRepository{}
		mockCSVRepo = &domainfakes.FakeCSVPowerConsumptionRepository{}
		mockPowerConsumptionService = NewPowerConsumptionService(mockMySQLRepo, mockCSVRepo)
	})

	Context("chekingKindPeriod", func() {
		It("should return monthly for 'monthly'", func() {
			result, err := mockPowerConsumptionService.ChekingKindPeriod("monthly")
			Expect(err).To(BeNil())
			Expect(result).To(Equal("monthly"))
		})

		It("should return weekly for 'weekly'", func() {
			result, err := mockPowerConsumptionService.ChekingKindPeriod("weekly")
			Expect(err).To(BeNil())
			Expect(result).To(Equal("weekly"))
		})

		It("should return daily for 'daily'", func() {
			result, err := mockPowerConsumptionService.ChekingKindPeriod("daily")
			Expect(err).To(BeNil())
			Expect(result).To(Equal("daily"))
		})

		It("should return error for invalid kind period", func() {
			result, err := mockPowerConsumptionService.ChekingKindPeriod("invalid")
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError(fmt.Errorf("Error: kind period not allowed %s", "invalid")))
			Expect(result).To(Equal(""))
		})
	})

})

var _ = Describe("PowerConsumptionService", func() {
	var (
		mockMySQLRepo               *domainfakes.FakeMySQLPowerConsumptionRepository
		mockCSVRepo                 *domainfakes.FakeCSVPowerConsumptionRepository
		mockPowerConsumptionService PowerConsumptionService
	)

	BeforeEach(func() {
		mockMySQLRepo = &domainfakes.FakeMySQLPowerConsumptionRepository{}
		mockCSVRepo = &domainfakes.FakeCSVPowerConsumptionRepository{}
		mockPowerConsumptionService = NewPowerConsumptionService(mockMySQLRepo, mockCSVRepo)
	})

	Context("ImportCsvToDatabase", func() {
		It("should successfully import CSV data to database", func() {
			csvData := []*domain.CSVUserConsumption{
				{
					ID:                 "1",
					MeterID:            "1",
					ActiveEnergy:       100.0,
					ReactiveEnergy:     50.0,
					CapacitiveReactive: 20.0,
					Solar:              30.0,
					Date:               "2023-08-01",
				},
				{
					ID:                 "2",
					MeterID:            "2",
					ActiveEnergy:       200.0,
					ReactiveEnergy:     60.0,
					CapacitiveReactive: 25.0,
					Solar:              35.0,
					Date:               "2023-08-02",
				},
			}

			mockCSVRepo.ConvertCSVToStructReturns(csvData, nil)

			mockMySQLRepo.CreatePowerConsumptionRecordsReturns(nil)

			err := mockPowerConsumptionService.ImportCsvToDatabase(nil)

			Expect(err).To(BeNil())
			Expect(mockCSVRepo.ConvertCSVToStructCallCount()).To(Equal(1))
			Expect(mockMySQLRepo.CreatePowerConsumptionRecordsCallCount()).To(Equal(1))
		})

		It("should return error when CSV conversion fails", func() {
			mockCSVRepo.ConvertCSVToStructReturns(nil, errors.New("Error reading CSV"))
			err := mockPowerConsumptionService.ImportCsvToDatabase(nil)

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("Error reading CSV"))

			Expect(mockMySQLRepo.CreatePowerConsumptionRecordsCallCount()).To(Equal(0))
		})

		It("should return error when CSV record conversion fails", func() {
			csvData := []*domain.CSVUserConsumption{
				{
					ID:                 "1",
					MeterID:            "1",
					ActiveEnergy:       100.0,
					ReactiveEnergy:     50.0,
					CapacitiveReactive: 20.0,
					Solar:              30.0,
					Date:               "2023-08-01",
				},
				{
					ID:                 "2",
					MeterID:            "2",
					ActiveEnergy:       200.0,
					ReactiveEnergy:     60.0,
					CapacitiveReactive: 25.0,
					Solar:              35.0,
					Date:               "2023-08-02",
				},
			}
			mockCSVRepo.ConvertCSVToStructReturns(csvData, nil)
			mockMySQLRepo.CreatePowerConsumptionRecordsReturns(errors.New("Error creating records"))
			err := mockPowerConsumptionService.ImportCsvToDatabase(nil)

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("Error creating records"))
			Expect(mockMySQLRepo.CreatePowerConsumptionRecordsCallCount()).To(Equal(1))
		})
	})

})

var _ = Describe("PowerConsumptionServiceImpl", func() {
	var (
		mockMySQLRepo *domainfakes.FakeMySQLPowerConsumptionRepository
		mockCSVRepo   *domainfakes.FakeCSVPowerConsumptionRepository
	)

	BeforeEach(func() {
		mockMySQLRepo = &domainfakes.FakeMySQLPowerConsumptionRepository{}
		mockCSVRepo = &domainfakes.FakeCSVPowerConsumptionRepository{}
	})

	Describe("GetConsumptionByMeterIDAndWindowTime", func() {
		const (
			meterIDs   = "1,2,3"
			startDate  = "2023-01-01"
			endDate    = "2023-01-31"
			kindPeriod = "monthly"
		)

		Context("when checking query parameters", func() {
			It("should return an error if checking query parameters fails", func() {
				expectedError := errors.New("checking query parameters failed")
				mockService := PowerConsumptionServiceImpl{
					mysqlRepository: mockMySQLRepo,
					csvRepository:   mockCSVRepo,
				}
				mockMySQLRepo.GetConsumptionByMeterIDAndWindowTimeReturns(nil, expectedError)

				result, err := mockService.GetConsumptionByMeterIDAndWindowTime(meterIDs, startDate, endDate, kindPeriod)
				Expect(result).To(BeNil())
				Expect(err).To(Equal(expectedError))
			})
		})

		Context("when getting consumption data from MySQL repository", func() {
			It("should return an error if getting data from MySQL repository fails", func() {
				expectedError := errors.New("getting data from MySQL repository failed")
				mockService := PowerConsumptionServiceImpl{
					mysqlRepository: mockMySQLRepo,
					csvRepository:   mockCSVRepo,
				}
				mockMySQLRepo.GetConsumptionByMeterIDAndWindowTimeReturns(nil, expectedError)
				mockMySQLRepo.GetConsumptionByMeterIDAndWindowTimeStub = func(startDate, endDate time.Time, meterID int) ([]domain.UserConsumption, error) {
					Expect(startDate).To(Equal(startDate))
					Expect(endDate).To(Equal(endDate))
					Expect(meterID).To(Equal(meterID))
					return nil, expectedError
				}

				result, err := mockService.GetConsumptionByMeterIDAndWindowTime(meterIDs, startDate, endDate, kindPeriod)
				Expect(result).To(BeNil())
				Expect(err).To(Equal(expectedError))
			})
		})

	})

})
