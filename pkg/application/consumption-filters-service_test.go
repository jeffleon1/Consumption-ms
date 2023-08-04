package application

import (
	"time"

	constants "github.com/jeffleon1/consumption-ms/internal/constans"
	"github.com/jeffleon1/consumption-ms/pkg/domain"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	shouldThirtyDays             string = "should correctly group a 30-day month"
	shouldSerializeDaysCorrectly string = "should serialize dates correctly"
	shouldThirtyOneDays          string = "should correctly group a 31-day month"
)

var _ = Describe("NewFilter Tests", func() {
	var (
		startDate time.Time
		endDate   time.Time
		data      []domain.UserConsumption
	)

	BeforeEach(func() {
		startDate = time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(2023, time.January, 31, 23, 59, 59, 999999999, time.UTC)
		data = []domain.UserConsumption{}
	})

	Context("when tipe is PeriodKindMonthly", func() {
		It("should return a MonthlyFilter", func() {
			filter := NewFilter(constants.PeriodKindMonthly, startDate, endDate, data)
			Expect(filter).To(BeAssignableToTypeOf(&MonthlyFilter{}))
		})
	})

	Context("when tipe is PeriodKindWeekly", func() {
		It("should return a WeeklyFilter", func() {
			filter := NewFilter(constants.PeriodKindWeekly, startDate, endDate, data)
			Expect(filter).To(BeAssignableToTypeOf(&WeeklyFilter{}))
		})
	})

	Context("when tipe is PeriodKindDaily", func() {
		It("should return a DailyFilter", func() {
			filter := NewFilter(constants.PeriodKindDaily, startDate, endDate, data)
			Expect(filter).To(BeAssignableToTypeOf(&DailyFilter{}))
		})
	})

	Context("when tipe is not valid", func() {
		It("should return nil", func() {
			filter := NewFilter("invalid_tipe", startDate, endDate, data)
			Expect(filter).To(BeNil())
		})
	})
})

var _ = Describe("days in month Tests", func() {
	var (
		filter *Filter
	)

	BeforeEach(func() {
		filter = &Filter{}
	})

	Context("daysInMonth", func() {
		It("should return the correct number of days for a given month and year", func() {
			month := 2
			year := 2023

			days := filter.daysInMonth(month, year)
			Expect(days).To(Equal(28))

			month = 12
			year = 2023

			days = filter.daysInMonth(month, year)
			Expect(days).To(Equal(31))

			month = 4
			year = 2023

			days = filter.daysInMonth(month, year)
			Expect(days).To(Equal(30))
		})
	})
})

var _ = Describe("DivideInformationByYears Tests", func() {
	var (
		filter *Filter
	)

	BeforeEach(func() {
		filter = &Filter{
			Data: []domain.UserConsumption{
				{Date: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				{Date: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)},
				{Date: time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC)},
				{Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
				{Date: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)},
			},
		}
	})

	Context("DivideInformationByYears", func() {
		It("should divide information by years and months", func() {
			result := filter.DivideInformationByYears()

			Expect(result).To(HaveLen(2))
			Expect(result[2022]).To(HaveKey(1))
			Expect(result[2022]).To(HaveKey(2))
			Expect(result[2022]).To(HaveKey(3))

			Expect(result[2023]).To(HaveKey(1))
			Expect(result[2023]).To(HaveKey(2))
		})

		It("should return an empty map if data is empty", func() {
			filter = &Filter{
				Data: []domain.UserConsumption{},
			}

			result := filter.DivideInformationByYears()

			Expect(result).To(BeEmpty())
		})

		It("should handle nil Data without errors", func() {
			filter = &Filter{
				Data: nil,
			}

			result := filter.DivideInformationByYears()

			Expect(result).To(BeEmpty())
		})

		It("should handle data with dates in reverse order", func() {
			filter = &Filter{
				Data: []domain.UserConsumption{
					{Date: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)},
					{Date: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				},
			}

			result := filter.DivideInformationByYears()

			Expect(result).To(HaveLen(2))

		})
	})
})

var _ = Describe("MatchConsumptionInTimeGroup Tests", func() {
	var (
		filter *Filter
	)

	BeforeEach(func() {
		filter = &Filter{}
	})

	Context("MatchConsumptionInTimeGroup", func() {
		It("should match consumption data correctly", func() {
			consumptions := []domain.UserConsumption{
				{Date: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				{Date: time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)},
				{Date: time.Date(2022, 1, 20, 0, 0, 0, 0, time.UTC)},
				{Date: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)},
				{Date: time.Date(2022, 2, 10, 0, 0, 0, 0, time.UTC)},
				{Date: time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC)},
			}

			timeGroups := []TimeGroupDivision{
				{
					InitDate:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
					FinishDate: time.Date(2022, 1, 20, 0, 0, 0, 0, time.UTC),
				},
				{
					InitDate:   time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
					FinishDate: time.Date(2022, 2, 28, 0, 0, 0, 0, time.UTC),
				},
			}

			result := filter.MatchConsumptionInTimeGroup(consumptions, timeGroups)
			Expect(result).To(HaveLen(2))

		})

		It("should return an empty slice if no matches are found", func() {
			consumptions := []domain.UserConsumption{
				{Date: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				{Date: time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)},
			}

			timeGroups := []TimeGroupDivision{
				{
					InitDate:   time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
					FinishDate: time.Date(2022, 2, 28, 0, 0, 0, 0, time.UTC),
				},
			}

			result := filter.MatchConsumptionInTimeGroup(consumptions, timeGroups)
			Expect(result).To(BeEmpty())
		})

		It("should handle nil consumptions without errors", func() {
			consumptions := []domain.UserConsumption(nil)
			timeGroups := []TimeGroupDivision{
				{
					InitDate:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
					FinishDate: time.Date(2022, 1, 20, 0, 0, 0, 0, time.UTC),
				},
			}

			result := filter.MatchConsumptionInTimeGroup(consumptions, timeGroups)

			Expect(result).To(BeEmpty())
		})

		It("should handle empty timeGroups without errors", func() {
			consumptions := []domain.UserConsumption{
				{Date: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				{Date: time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)},
			}

			timeGroups := []TimeGroupDivision(nil)

			result := filter.MatchConsumptionInTimeGroup(consumptions, timeGroups)

			Expect(result).To(BeEmpty())
		})
	})
})

var _ = Describe("ReduceInformation Tests", func() {
	var (
		filter *Filter
	)

	BeforeEach(func() {
		filter = &Filter{}
	})

	Context("ReduceInformation", func() {
		It("should reduce consumption data correctly", func() {
			weeklyGroupConsumptions := []*ConsumptionEnergy{
				{
					ActiveEnergy:       100,
					ReactiveEnergy:     50,
					CapacitiveReactive: 20,
					Exported:           30,
					Data: []domain.UserConsumption{
						{
							ActiveEnergy:       20,
							ReactiveEnergy:     10,
							CapacitiveReactive: 5,
							Solar:              15,
						},
						{
							ActiveEnergy:       30,
							ReactiveEnergy:     15,
							CapacitiveReactive: 7,
							Solar:              20,
						},
					},
				},
				{
					ActiveEnergy:       200,
					ReactiveEnergy:     100,
					CapacitiveReactive: 40,
					Exported:           60,
					Data: []domain.UserConsumption{
						{
							ActiveEnergy:       40,
							ReactiveEnergy:     20,
							CapacitiveReactive: 10,
							Solar:              30,
						},
					},
				},
			}

			filter.ReduceInformation(weeklyGroupConsumptions)

			Expect(weeklyGroupConsumptions[0].ActiveEnergy).To(Equal(float64(150)))      // 100 + 20 + 30
			Expect(weeklyGroupConsumptions[0].ReactiveEnergy).To(Equal(float64(75)))     // 50 + 10 + 15
			Expect(weeklyGroupConsumptions[0].CapacitiveReactive).To(Equal(float64(32))) // 20 + 5 + 7
			Expect(weeklyGroupConsumptions[0].Exported).To(Equal(float64(65)))           // 30 + 15 + 20

			Expect(weeklyGroupConsumptions[1].ActiveEnergy).To(Equal(float64(240)))      // 200 + 40
			Expect(weeklyGroupConsumptions[1].ReactiveEnergy).To(Equal(float64(120)))    // 100 + 20
			Expect(weeklyGroupConsumptions[1].CapacitiveReactive).To(Equal(float64(50))) // 40 + 10
			Expect(weeklyGroupConsumptions[1].Exported).To(Equal(float64(90)))           // 60 + 30
		})

		It("should handle empty weeklyGroupConsumptions without errors", func() {
			weeklyGroupConsumptions := []*ConsumptionEnergy{}
			filter.ReduceInformation(weeklyGroupConsumptions)
		})

		It("should handle nil weeklyGroupConsumptions without errors", func() {
			weeklyGroupConsumptions := []*ConsumptionEnergy(nil)
			filter.ReduceInformation(weeklyGroupConsumptions)
		})
	})
})

var _ = Describe("Monthly Group division", func() {
	var (
		monthlyFilter *MonthlyFilter
	)

	BeforeEach(func() {
		monthlyFilter = &MonthlyFilter{}
	})

	Context("GroupDivision", func() {
		It("should correctly group a 30-day month", func() {
			month := 4
			year := 2023
			result := monthlyFilter.GroupDivision(month, year)

			Expect(result).To(HaveLen(1))
			initDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			finishDate := time.Date(year, time.Month(month), 30, 23, 59, 59, 59, time.UTC)
			Expect(result[0].InitDate).To(Equal(initDate))
			Expect(result[0].FinishDate).To(Equal(finishDate))
		})

		It("should correctly group a 31-day month", func() {
			month := 5
			year := 2023

			result := monthlyFilter.GroupDivision(month, year)
			Expect(result).To(HaveLen(1))

			initDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			finishDate := time.Date(year, time.Month(month), 31, 23, 59, 59, 59, time.UTC)
			Expect(result[0].InitDate).To(Equal(initDate))
			Expect(result[0].FinishDate).To(Equal(finishDate))
		})

		It("should correctly group a February of a leap year", func() {
			month := 2
			year := 2024

			result := monthlyFilter.GroupDivision(month, year)
			Expect(result).To(HaveLen(1))
			initDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			finishDate := time.Date(year, time.Month(month), 29, 23, 59, 59, 59, time.UTC)
			Expect(result[0].InitDate).To(Equal(initDate))
			Expect(result[0].FinishDate).To(Equal(finishDate))
		})

		It("should correctly handle an invalid month", func() {
			month := 13
			year := 2023

			result := monthlyFilter.GroupDivision(month, year)
			Expect(result).To(BeEmpty())
		})

		Context("GroupsSerializedToString", func() {
			It("should serialize dates correctly", func() {
				startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2023, 1, 31, 23, 59, 59, 59, time.UTC)
				result := monthlyFilter.GroupsSerializedToString(startDate, endDate)
				expectedResult := "Jan 2023"
				Expect(result).To(Equal(expectedResult))
			})

			It("should handle different years correctly", func() {
				startDate := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2024, 2, 29, 23, 59, 59, 59, time.UTC)
				result := monthlyFilter.GroupsSerializedToString(startDate, endDate)

				expectedResult := "Feb 2024"
				Expect(result).To(Equal(expectedResult))
			})

			It("should handle endDate in a different year correctly", func() {
				startDate := time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2023, 1, 31, 23, 59, 59, 59, time.UTC)
				result := monthlyFilter.GroupsSerializedToString(startDate, endDate)

				expectedResult := "Dec 2022"
				Expect(result).To(Equal(expectedResult))
			})

			It("should handle an invalid date range correctly", func() {
				startDate := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2023, 1, 1, 23, 59, 59, 59, time.UTC)

				result := monthlyFilter.GroupsSerializedToString(startDate, endDate)
				Expect(result).To(BeEmpty())
			})
		})
	})
})

var _ = Describe("Daily filter", func() {
	var (
		dailyFilter *DailyFilter
	)

	BeforeEach(func() {
		dailyFilter = &DailyFilter{}
	})

	Context("GroupDivision", func() {
		It(shouldThirtyDays, func() {
			month := 4
			year := 2023

			result := dailyFilter.GroupDivision(month, year)
			Expect(result).To(HaveLen(30))
			initDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			for i, group := range result {
				Expect(group.InitDate).To(Equal(initDate.AddDate(0, 0, i)))
				Expect(group.FinishDate).To(Equal(initDate.AddDate(0, 0, i).AddDate(0, 0, 1).Add(-time.Second)))
			}
		})

		It(shouldThirtyOneDays, func() {
			month := 5
			year := 2023

			result := dailyFilter.GroupDivision(month, year)
			Expect(result).To(HaveLen(31))
			initDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			for i, group := range result {
				Expect(group.InitDate).To(Equal(initDate.AddDate(0, 0, i)))
				Expect(group.FinishDate).To(Equal(initDate.AddDate(0, 0, i).AddDate(0, 0, 1).Add(-time.Second)))
			}
		})

		It("should handle February of a leap year correctly", func() {
			month := 2
			year := 2024

			result := dailyFilter.GroupDivision(month, year)
			Expect(result).To(HaveLen(29))
			initDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			for i, group := range result {
				Expect(group.InitDate).To(Equal(initDate.AddDate(0, 0, i)))
				Expect(group.FinishDate).To(Equal(initDate.AddDate(0, 0, i).AddDate(0, 0, 1).Add(-time.Second)))
			}
		})

		It("should handle an invalid month correctly", func() {
			month := 13
			year := 2023
			result := dailyFilter.GroupDivision(month, year)
			Expect(result).To(BeEmpty())
		})
	})

	Context("GroupsSerializedToString", func() {
		It(shouldSerializeDaysCorrectly, func() {
			startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(2023, 1, 31, 23, 59, 59, 59, time.UTC)
			result := dailyFilter.GroupsSerializedToString(startDate, endDate)
			expectedResult := "Jan 1"
			Expect(result).To(Equal(expectedResult))
		})

		It("should handle endDate before startDate correctly", func() {
			startDate := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(2023, 1, 1, 23, 59, 59, 59, time.UTC)
			result := dailyFilter.GroupsSerializedToString(startDate, endDate)

			Expect(result).To(BeEmpty())
		})
	})
})

var _ = Describe("Weekly filter", func() {
	var (
		weeklyFilter *WeeklyFilter
	)

	BeforeEach(func() {
		weeklyFilter = &WeeklyFilter{}
	})

	Context("GroupDivision", func() {
		It(shouldThirtyDays, func() {
			month := 4
			year := 2023
			result := weeklyFilter.GroupDivision(month, year)
			Expect(result).To(HaveLen(5))
		})

		It(shouldThirtyOneDays, func() {
			month := 5
			year := 2023

			result := weeklyFilter.GroupDivision(month, year)
			Expect(result).To(HaveLen(5))
		})

		It("should handle February of a leap year correctly", func() {
			month := 2
			year := 2024
			result := weeklyFilter.GroupDivision(month, year)
			Expect(result).To(HaveLen(5))
		})

		It("should handle an invalid month correctly", func() {
			month := 13
			year := 2023
			result := weeklyFilter.GroupDivision(month, year)
			Expect(result).To(BeEmpty())
		})

		Context("GroupsSerializedToString", func() {
			It(shouldSerializeDaysCorrectly, func() {
				startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2023, 1, 7, 23, 59, 59, 59, time.UTC)
				result := weeklyFilter.GroupsSerializedToString(startDate, endDate)
				expectedResult := "Jan 1 - Jan 7"
				Expect(result).To(Equal(expectedResult))
			})

			It("should handle endDate before startDate correctly", func() {
				startDate := time.Date(2023, 1, 7, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2023, 1, 1, 23, 59, 59, 59, time.UTC)
				result := weeklyFilter.GroupsSerializedToString(startDate, endDate)
				Expect(result).To(BeEmpty())
			})
		})
	})
})
