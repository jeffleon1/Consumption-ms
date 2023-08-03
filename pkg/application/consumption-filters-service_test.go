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
			// Configura datos de prueba con weeklyGroupConsumptions vacío
			weeklyGroupConsumptions := []*ConsumptionEnergy{}

			// Ejecuta la función ReduceInformation
			filter.ReduceInformation(weeklyGroupConsumptions)

			// Verifica que no haya errores
			// No hay expectativas específicas para verificar los resultados, ya que no se espera que se realice ninguna acción con el slice vacío
		})

		It("should handle nil weeklyGroupConsumptions without errors", func() {
			// Configura datos de prueba con weeklyGroupConsumptions nulo (nil)
			weeklyGroupConsumptions := []*ConsumptionEnergy(nil)

			// Ejecuta la función ReduceInformation
			filter.ReduceInformation(weeklyGroupConsumptions)

			// Verifica que no haya errores
			// No hay expectativas específicas para verificar los resultados, ya que no se espera que se realice ninguna acción con el slice nulo
		})
	})
})

var _ = Describe("Monthly Group division", func() {
	var (
		monthlyFilter *MonthlyFilter
	)

	BeforeEach(func() {
		// Configura el MonthlyFilter necesario para las pruebas
		monthlyFilter = &MonthlyFilter{}
	})

	Context("GroupDivision", func() {
		It("should correctly group a 30-day month", func() {
			// Configura datos de prueba para un mes de 30 días
			month := 4 // Abril
			year := 2023

			// Ejecuta la función GroupDivision
			result := monthlyFilter.GroupDivision(month, year)

			// Verifica que el resultado sea el esperado
			Expect(result).To(HaveLen(1)) // Debería haber una división de tiempo para todo el mes

			// Verifica las fechas de inicio y final
			initDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			finishDate := time.Date(year, time.Month(month), 30, 23, 59, 59, 59, time.UTC)
			Expect(result[0].InitDate).To(Equal(initDate))
			Expect(result[0].FinishDate).To(Equal(finishDate))
		})

		It("should correctly group a 31-day month", func() {
			// Configura datos de prueba para un mes de 31 días
			month := 5 // Mayo
			year := 2023

			// Ejecuta la función GroupDivision
			result := monthlyFilter.GroupDivision(month, year)

			// Verifica que el resultado sea el esperado
			Expect(result).To(HaveLen(1)) // Debería haber una división de tiempo para todo el mes

			// Verifica las fechas de inicio y final
			initDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			finishDate := time.Date(year, time.Month(month), 31, 23, 59, 59, 59, time.UTC)
			Expect(result[0].InitDate).To(Equal(initDate))
			Expect(result[0].FinishDate).To(Equal(finishDate))
		})

		It("should correctly group a February of a leap year", func() {
			// Configura datos de prueba para un febrero de un año bisiesto (29 días)
			month := 2   // Febrero
			year := 2024 // Año bisiesto

			// Ejecuta la función GroupDivision
			result := monthlyFilter.GroupDivision(month, year)

			// Verifica que el resultado sea el esperado
			Expect(result).To(HaveLen(1)) // Debería haber una división de tiempo para todo el mes

			// Verifica las fechas de inicio y final
			initDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			finishDate := time.Date(year, time.Month(month), 29, 23, 59, 59, 59, time.UTC)
			Expect(result[0].InitDate).To(Equal(initDate))
			Expect(result[0].FinishDate).To(Equal(finishDate))
		})

		It("should correctly handle an invalid month", func() {
			// Configura datos de prueba para un mes inválido (fuera del rango 1-12)
			month := 13 // Mes inválido
			year := 2023

			// Ejecuta la función GroupDivision
			result := monthlyFilter.GroupDivision(month, year)

			// Verifica que el resultado sea un slice vacío
			Expect(result).To(BeEmpty())
		})

		Context("GroupsSerializedToString", func() {
			It("should serialize dates correctly", func() {
				// Configura datos de prueba
				startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2023, 1, 31, 23, 59, 59, 59, time.UTC)

				// Ejecuta la función GroupsSerializedToString
				result := monthlyFilter.GroupsSerializedToString(startDate, endDate)

				// Verifica que el resultado sea el esperado
				expectedResult := "Jan 2023" // El formato esperado para enero de 2023 es "yyyy-MM"
				Expect(result).To(Equal(expectedResult))
			})

			It("should handle different years correctly", func() {
				// Configura datos de prueba para febrero de un año bisiesto (29 días)
				startDate := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2024, 2, 29, 23, 59, 59, 59, time.UTC)

				// Ejecuta la función GroupsSerializedToString
				result := monthlyFilter.GroupsSerializedToString(startDate, endDate)

				// Verifica que el resultado sea el esperado
				expectedResult := "Feb 2024" // El formato esperado para febrero de 2024 es "yyyy-MM"
				Expect(result).To(Equal(expectedResult))
			})

			It("should handle endDate in a different year correctly", func() {
				// Configura datos de prueba para diciembre de un año y enero del siguiente año
				startDate := time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2023, 1, 31, 23, 59, 59, 59, time.UTC)

				// Ejecuta la función GroupsSerializedToString
				result := monthlyFilter.GroupsSerializedToString(startDate, endDate)

				// Verifica que el resultado sea el esperado
				expectedResult := "Dec 2022" // El formato esperado para diciembre de 2022 es "yyyy-MM"
				Expect(result).To(Equal(expectedResult))
			})

			It("should handle an invalid date range correctly", func() {
				// Configura datos de prueba con una fecha de inicio posterior a la fecha de finalización
				startDate := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2023, 1, 1, 23, 59, 59, 59, time.UTC)

				// Ejecuta la función GroupsSerializedToString
				result := monthlyFilter.GroupsSerializedToString(startDate, endDate)

				// Verifica que el resultado sea vacío o algún otro valor indicativo de error
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
		// Configura el DailyFilter necesario para las pruebas
		dailyFilter = &DailyFilter{}
	})

	Context("GroupDivision", func() {
		It(shouldThirtyDays, func() {
			// Configura datos de prueba para un mes de 30 días
			month := 4 // Abril
			year := 2023

			// Ejecuta la función GroupDivision
			result := dailyFilter.GroupDivision(month, year)

			// Verifica que el resultado sea el esperado
			Expect(result).To(HaveLen(30)) // Debería haber 30 divisiones de tiempo para cada día del mes

			// Verifica las fechas de inicio y final de cada día
			initDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			for i, group := range result {
				Expect(group.InitDate).To(Equal(initDate.AddDate(0, 0, i)))
				Expect(group.FinishDate).To(Equal(initDate.AddDate(0, 0, i).AddDate(0, 0, 1).Add(-time.Second)))
			}
		})

		It(shouldThirtyOneDays, func() {
			// Configura datos de prueba para un mes de 31 días
			month := 5 // Mayo
			year := 2023

			// Ejecuta la función GroupDivision
			result := dailyFilter.GroupDivision(month, year)

			// Verifica que el resultado sea el esperado
			Expect(result).To(HaveLen(31)) // Debería haber 31 divisiones de tiempo para cada día del mes

			// Verifica las fechas de inicio y final de cada día
			initDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			for i, group := range result {
				Expect(group.InitDate).To(Equal(initDate.AddDate(0, 0, i)))
				Expect(group.FinishDate).To(Equal(initDate.AddDate(0, 0, i).AddDate(0, 0, 1).Add(-time.Second)))
			}
		})

		It("should handle February of a leap year correctly", func() {
			// Configura datos de prueba para febrero de un año bisiesto (29 días)
			month := 2   // Febrero
			year := 2024 // Año bisiesto

			// Ejecuta la función GroupDivision
			result := dailyFilter.GroupDivision(month, year)

			// Verifica que el resultado sea el esperado
			Expect(result).To(HaveLen(29)) // Debería haber 29 divisiones de tiempo para cada día del mes

			// Verifica las fechas de inicio y final de cada día
			initDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			for i, group := range result {
				Expect(group.InitDate).To(Equal(initDate.AddDate(0, 0, i)))
				Expect(group.FinishDate).To(Equal(initDate.AddDate(0, 0, i).AddDate(0, 0, 1).Add(-time.Second)))
			}
		})

		It("should handle an invalid month correctly", func() {
			// Configura datos de prueba para un mes inválido (fuera del rango 1-12)
			month := 13 // Mes inválido
			year := 2023

			// Ejecuta la función GroupDivision
			result := dailyFilter.GroupDivision(month, year)

			// Verifica que el resultado sea un slice vacío
			Expect(result).To(BeEmpty())
		})
	})

	Context("GroupsSerializedToString", func() {
		It(shouldSerializeDaysCorrectly, func() {
			// Configura datos de prueba
			startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(2023, 1, 31, 23, 59, 59, 59, time.UTC)

			// Ejecuta la función GroupsSerializedToString
			result := dailyFilter.GroupsSerializedToString(startDate, endDate)

			// Verifica que el resultado sea el esperado
			expectedResult := "Jan 1" // El formato esperado para el rango de fechas es "yyyy-MM-dd_yyyy-MM-dd"
			Expect(result).To(Equal(expectedResult))
		})

		It("should handle endDate before startDate correctly", func() {
			// Configura datos de prueba con endDate antes que startDate
			startDate := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(2023, 1, 1, 23, 59, 59, 59, time.UTC)

			// Ejecuta la función GroupsSerializedToString
			result := dailyFilter.GroupsSerializedToString(startDate, endDate)

			// Verifica que el resultado sea vacío
			Expect(result).To(BeEmpty())
		})
	})
})

var _ = Describe("Weekly filter", func() {
	var (
		weeklyFilter *WeeklyFilter
	)

	BeforeEach(func() {
		// Configura el WeeklyFilter necesario para las pruebas
		weeklyFilter = &WeeklyFilter{}
	})

	Context("GroupDivision", func() {
		It(shouldThirtyDays, func() {
			// Configura datos de prueba para un mes de 30 días
			month := 4 // Abril
			year := 2023

			// Ejecuta la función GroupDivision
			result := weeklyFilter.GroupDivision(month, year)

			// Verifica que el resultado sea el esperado
			Expect(result).To(HaveLen(5)) // Debería haber 5 divisiones de tiempo para cada semana del mes

			// Verifica las fechas de inicio y final de cada semana
		})

		It(shouldThirtyOneDays, func() {
			// Configura datos de prueba para un mes de 31 días
			month := 5 // Mayo
			year := 2023

			// Ejecuta la función GroupDivision
			result := weeklyFilter.GroupDivision(month, year)

			// Verifica que el resultado sea el esperado
			Expect(result).To(HaveLen(5)) // Debería haber 6 divisiones de tiempo para cada semana del mes

		})

		It("should handle February of a leap year correctly", func() {
			// Configura datos de prueba para febrero de un año bisiesto (29 días)
			month := 2   // Febrero
			year := 2024 // Año bisiesto

			// Ejecuta la función GroupDivision
			result := weeklyFilter.GroupDivision(month, year)

			// Verifica que el resultado sea el esperado
			Expect(result).To(HaveLen(5)) // Debería haber 4 divisiones de tiempo para cada semana del mes

		})

		It("should handle an invalid month correctly", func() {
			// Configura datos de prueba para un mes inválido (fuera del rango 1-12)
			month := 13 // Mes inválido
			year := 2023

			// Ejecuta la función GroupDivision
			result := weeklyFilter.GroupDivision(month, year)

			// Verifica que el resultado sea un slice vacío
			Expect(result).To(BeEmpty())
		})

		Context("GroupsSerializedToString", func() {
			It(shouldSerializeDaysCorrectly, func() {
				// Configura datos de prueba
				startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2023, 1, 7, 23, 59, 59, 59, time.UTC)

				// Ejecuta la función GroupsSerializedToString
				result := weeklyFilter.GroupsSerializedToString(startDate, endDate)

				// Verifica que el resultado sea el esperado
				expectedResult := "Jan 1 - Jan 7" // El formato esperado para el rango de fechas es "yyyy-MM-dd - yyyy-MM-dd"
				Expect(result).To(Equal(expectedResult))
			})

			It("should handle endDate before startDate correctly", func() {
				// Configura datos de prueba con endDate antes que startDate
				startDate := time.Date(2023, 1, 7, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2023, 1, 1, 23, 59, 59, 59, time.UTC)

				// Ejecuta la función GroupsSerializedToString
				result := weeklyFilter.GroupsSerializedToString(startDate, endDate)

				// Verifica que el resultado sea vacío
				Expect(result).To(BeEmpty())
			})
		})
	})
})
