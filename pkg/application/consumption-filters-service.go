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

// NewFilter: Factory to create filterss
//
// Parámeters:
// tipe - type of filter
// startDate - start date to retrieve or filter the information
// endDate - type of filter
// tipe - type of filter
//
// Returns:
// The struct that repesents the database domain
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

// GetConsumptionData: this function is the main function in this file becuase have all the logic to retrieve
// the records then organize that records and return the information in the way that we want
//
// Parámeters:
// filter - is an interface that allow do the process with all types of filters no matter what kind of filter is
//
// Returns:
// The Serializer that is a kind of structure that has all the attributes that we need to serialize in consumption serializer
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

// daysInMonth: get and month and year and return the number of days for this especific month in this specific year
//
// Parámeters:
// month - month to retrieve the days
// year - especify the year and month to retrieve the days
//
// Returns:
// return the days in this especific month in this especific year
func (f *Filter) daysInMonth(month, year int) int {
	firstDayOfNextMonth := time.Date(year, time.Month(month)+1, 1, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := firstDayOfNextMonth.Add(-time.Second)
	return lastDayOfMonth.Day()
}

// DivideInformationByYears: divide the information retrieve in the database and divided
// in a map by year --> month ---> records in month
//
// Returns:
// return a map with all the information by year --> month ---> records in month
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

// MatchConsumptionInTimeGroup: do the match between the userconsumption and
// the group division no matter if it's a group division by monthly, weekly or daily
//
// Parameters:
// consumptions: has the consumption information
// timeGroups: has the time groups by monthly weekly of daily depends of the type of filter
//
// Returns:
// return the information matched between timegroups and userconsumption
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

// ReduceInformation: reduce the information a only one record by group division
//
// Parameters:
// groupConsumptions: has the information matched between groups and information
//
// Returns:
// return reduced and one record by group division
func (f *Filter) ReduceInformation(groupConsumptions []*ConsumptionEnergy) {
	for _, groupConsumption := range groupConsumptions {
		for _, group := range groupConsumption.Data {
			groupConsumption.ActiveEnergy += group.ActiveEnergy
			groupConsumption.ReactiveEnergy += group.ReactiveEnergy
			groupConsumption.CapacitiveReactive += group.CapacitiveReactive
			groupConsumption.Exported += group.Solar
		}
	}
	logrus.Info("Reduce information is done")
}

// Monthly filter
// GroupDivision: do the group division for a monthly filter
//
// Parameters:
// month
// year
//
// Returns:
// return the time group division for a monthly filter
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
// GroupDivision: do the group division for a daily filter
//
// Parameters:
// month
// year
//
// Returns:
// return the time group division for a daily filter
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

// GroupsSerializedToString: serialize the date in a way that need the filter example "2023-01-02" --> "Jan 2"
//
// Parameters:
// startDate
// endDate
//
// Returns:
// return the date with in a correct way "Jan 2"
func (d *DailyFilter) GroupsSerializedToString(startDate time.Time, endDate time.Time) string {
	if startDate.After(endDate) {
		return ""
	}
	return domain.TimeTostr(startDate, constants.DateFormatWeeklyAndDailyPeriod)
}

// Week Filter
// GroupDivision: do the group division for a weekly filter
//
// Parameters:
// month
// year
//
// Returns:
// return the time group division for a weekly filter
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

// GroupsSerializedToString: serialize the date in a way that need the filter example "2023-01-02" "2023-01-08" --> "Jan 2 - Jan 8"
//
// Parameters:
// startDate
// endDate
//
// Returns:
// return the date with in a correct way "Jan 2 - Jan 8"
func (w *WeeklyFilter) GroupsSerializedToString(startDate time.Time, endDate time.Time) string {
	if startDate.After(endDate) {
		return ""
	}
	startDateString := domain.TimeTostr(startDate, constants.DateFormatWeeklyAndDailyPeriod)
	endDateString := domain.TimeTostr(endDate, constants.DateFormatWeeklyAndDailyPeriod)
	return fmt.Sprintf("%s - %s", startDateString, endDateString)
}
