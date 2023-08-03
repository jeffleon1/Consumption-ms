package application

import (
	"fmt"
	"sort"
	"time"

	constants "github.com/jeffleon1/consumption-ms/internal/constans"
	"github.com/jeffleon1/consumption-ms/pkg/domain"
	"github.com/sirupsen/logrus"
)

type FilterOperations interface {
	FilterCommonOperations
	FilterAdditionalOperations
}

type FilterCommonOperations interface {
	DivideInformationByYears() map[int]map[int][]domain.UserConsumption
	daysInMonth(month, year int) int
	MatchConsumptionInTimeGroup(consumptions []domain.UserConsumption, timeGroups []TimeGroupDivision) []*ConsumptionEnergy
	ReduceInformation(weeklyGroupConsumptions []*ConsumptionEnergy)
}

type FilterAdditionalOperations interface {
	GroupDivision(month, year int) []TimeGroupDivision
	GroupsSerializedToString(time.Time, time.Time) string
}

type Filter struct {
	StartDate time.Time
	EndDate   time.Time
	Data      []domain.UserConsumption
}

type TimeGroupDivision struct {
	InitDate   time.Time
	FinishDate time.Time
}

type ConsumptionEnergy struct {
	StartDate          time.Time
	EndDate            time.Time
	Data               []domain.UserConsumption
	ActiveEnergy       float64
	ReactiveEnergy     float64
	CapacitiveReactive float64
	Exported           float64
}

type Serializer struct {
	Period             []string  `json:"period"`
	MeterID            int       `json:"meter_id"`
	Active             []float64 `json:"active"`
	ReactiveInductive  []float64 `json:"reactive_inductive"`
	ReactiveCapacitive []float64 `json:"reactive_capacitive"`
	Exported           []float64 `json:"exported"`
}

type MonthlyFilter struct {
	Filter
}

type DailyFilter struct {
	Filter
}

type WeeklyFilter struct {
	Filter
}

func NewFilter(tipe string, startDate, endDate time.Time, data []domain.UserConsumption) FilterOperations {
	switch tipe {
	case constants.PeriodKindMonthly:
		return &MonthlyFilter{Filter{StartDate: startDate, EndDate: endDate, Data: data}}
	case constants.PeriodKindWeekly:
		return &WeeklyFilter{Filter{StartDate: startDate, EndDate: endDate, Data: data}}
	case constants.PeriodKindDaily:
		return &DailyFilter{Filter{StartDate: startDate, EndDate: endDate, Data: data}}
	default:
		return nil
	}
}

func GetConsumptionData(filter FilterOperations) Serializer {
	consumptionByYear := filter.DivideInformationByYears()
	var consumptionEnergy []*ConsumptionEnergy
	var objectSerializer Serializer
	for year, consumptionYear := range consumptionByYear {
		for month, conconsumptionInMonth := range consumptionYear {
			dailyGroups := filter.GroupDivision(month, year)
			informationMatched := filter.MatchConsumptionInTimeGroup(conconsumptionInMonth, dailyGroups)
			filter.ReduceInformation(informationMatched)
			if informationMatched != nil {
				consumptionEnergy = append(consumptionEnergy, informationMatched...)
			}
		}
	}

	for _, serializer := range consumptionEnergy {
		periodString := filter.GroupsSerializedToString(serializer.StartDate, serializer.EndDate)
		objectSerializer.Period = append(objectSerializer.Period, periodString)
		objectSerializer.Active = append(objectSerializer.Active, serializer.ActiveEnergy)
		objectSerializer.Exported = append(objectSerializer.Exported, serializer.Exported)
		objectSerializer.ReactiveInductive = append(objectSerializer.ReactiveInductive, serializer.ReactiveEnergy)
		objectSerializer.ReactiveCapacitive = append(objectSerializer.ReactiveCapacitive, serializer.CapacitiveReactive)
	}
	return objectSerializer
}

// Common Filter Operations
func (f *Filter) daysInMonth(month, year int) int {
	firstDayOfNextMonth := time.Date(year, time.Month(month)+1, 1, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := firstDayOfNextMonth.Add(-time.Second)
	return lastDayOfMonth.Day()
}

func (f *Filter) DivideInformationByYears() map[int]map[int][]domain.UserConsumption {
	data := f.Data
	if len(data) == 0 {
		return make(map[int]map[int][]domain.UserConsumption)
	}
	objectYearInformation := make(map[int]map[int][]domain.UserConsumption)
	objectMonthInformation := make(map[int][]domain.UserConsumption)
	sort.Slice(data, func(i, j int) bool {
		return data[i].Date.Before(data[j].Date)
	})
	month := data[0].Date.Month()
	year := data[0].Date.Year()
	for _, value := range data {
		if month != value.Date.Month() {
			month = value.Date.Month()
		}
		objectMonthInformation[int(month)] = append(objectMonthInformation[int(month)], value)
	}

	for _, value := range data {
		if year != value.Date.Year() {
			year = value.Date.Year()
		}
		objectYearInformation[int(year)] = objectMonthInformation
	}

	return objectYearInformation
}

func (f *Filter) MatchConsumptionInTimeGroup(consumptions []domain.UserConsumption, timeGroups []TimeGroupDivision) []*ConsumptionEnergy {
	var serializer []*ConsumptionEnergy
	for _, timeGroup := range timeGroups {
		var data []domain.UserConsumption
		for _, objectConsumption := range consumptions {
			if objectConsumption.Date.After(timeGroup.InitDate) && objectConsumption.Date.Before(timeGroup.FinishDate) {
				data = append(data, objectConsumption)
			}
		}
		if len(data) > 0 {
			serializer = append(serializer, &ConsumptionEnergy{
				StartDate: timeGroup.InitDate,
				EndDate:   timeGroup.FinishDate,
				Data:      data,
			})
		}
	}

	logrus.Info("Match consumption is done")
	return serializer
}

func (f *Filter) ReduceInformation(weeklyGroupConsumptions []*ConsumptionEnergy) {
	for _, weeklyGroupConsumption := range weeklyGroupConsumptions {
		for _, weeklyGroup := range weeklyGroupConsumption.Data {
			weeklyGroupConsumption.ActiveEnergy += weeklyGroup.ActiveEnergy
			weeklyGroupConsumption.ReactiveEnergy += weeklyGroup.ReactiveEnergy
			weeklyGroupConsumption.CapacitiveReactive += weeklyGroup.CapacitiveReactive
			weeklyGroupConsumption.Exported += weeklyGroup.Solar
		}
	}
	logrus.Info("Reduce information is done")
}

// Monthly filter
func (m *MonthlyFilter) GroupDivision(month, year int) []TimeGroupDivision {
	daysInMonth := m.daysInMonth(month, year)
	var monthGroups []TimeGroupDivision
	if month < 1 || month > 12 {
		return []TimeGroupDivision{}
	}
	initialDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDate := time.Date(year, time.Month(month), daysInMonth, 23, 59, 59, 59, time.UTC)
	monthGroups = append(monthGroups, TimeGroupDivision{
		InitDate:   initialDate,
		FinishDate: lastDate,
	})
	return monthGroups
}

func (m *MonthlyFilter) GroupsSerializedToString(startDate time.Time, endDate time.Time) string {
	if startDate.After(endDate) {
		return ""
	}
	return domain.TimeTostr(startDate, constants.DateFormatMonthlyPeriod)
}

// Dayly filter
func (d *DailyFilter) GroupDivision(month, year int) []TimeGroupDivision {
	daysInMonth := d.daysInMonth(month, year)
	var dayGroups []TimeGroupDivision
	if month < 1 || month > 12 {
		return []TimeGroupDivision{}
	}
	for i := 1; i <= daysInMonth; i++ {
		initialDay := i
		initialDate := time.Date(year, time.Month(month), initialDay, 0, 0, 0, 0, time.UTC)
		lastDate := initialDate.AddDate(0, 0, 1).Add(-time.Second)
		dayGroups = append(dayGroups, TimeGroupDivision{
			InitDate:   initialDate,
			FinishDate: lastDate,
		})
	}
	logrus.Info("day group division is done")
	return dayGroups
}

func (d *DailyFilter) GroupsSerializedToString(startDate time.Time, endDate time.Time) string {
	if startDate.After(endDate) {
		return ""
	}
	return domain.TimeTostr(startDate, constants.DateFormatWeeklyAndDailyPeriod)
}

// Week Filter
func (w *WeeklyFilter) GroupDivision(month, year int) []TimeGroupDivision {
	daysInMonth := w.daysInMonth(month, year)
	var weekGroups []TimeGroupDivision
	if month < 1 || month > 12 {
		return []TimeGroupDivision{}
	}
	for i := 1; i <= daysInMonth; i += 7 {
		initialDay := i
		initialDate := time.Date(year, time.Month(month), initialDay, 0, 0, 0, 0, time.UTC)
		lastDate := initialDate.AddDate(0, 0, 7).Add(-time.Second)
		if initialDay+6 > daysInMonth {
			lastDate = time.Date(year, time.Month(month), daysInMonth, 23, 59, 59, 59, time.UTC)
		}
		weekGroups = append(weekGroups, TimeGroupDivision{
			InitDate:   initialDate,
			FinishDate: lastDate,
		})
	}
	return weekGroups
}

func (w *WeeklyFilter) GroupsSerializedToString(startDate time.Time, endDate time.Time) string {
	if startDate.After(endDate) {
		return ""
	}
	startDateString := domain.TimeTostr(startDate, constants.DateFormatWeeklyAndDailyPeriod)
	endDateString := domain.TimeTostr(endDate, constants.DateFormatWeeklyAndDailyPeriod)
	return fmt.Sprintf("%s - %s", startDateString, endDateString)
}
