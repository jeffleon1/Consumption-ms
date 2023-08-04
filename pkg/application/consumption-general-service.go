package application

import (
	"fmt"
	"mime/multipart"
	"strings"
	"sync"
	"time"

	constants "github.com/jeffleon1/consumption-ms/internal/constans"
	"github.com/jeffleon1/consumption-ms/pkg/domain"
	"github.com/sirupsen/logrus"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . PowerConsumptionService
type PowerConsumptionService interface {
	GetConsumptionByMeterIDAndWindowTime(meterIDs string, kindPeriod string, startDate string, endDate string) ([]Serializer, error)
	ImportCsvToDatabase(file *multipart.File) error
	ChekingKindPeriod(kindPeriod string) (string, error)
	CheckingQueryParamConstrains(meterIDs string, kindPeriod string, startDate string, endDate string) (*domain.UserConsumptionQueryParams, error)
}

type PowerConsumptionServiceImpl struct {
	mysqlRepository domain.MySQLPowerConsumptionRepository
	csvRepository   domain.CSVPowerConsumptionRepository
}

func NewPowerConsumptionService(mysqlRepository domain.MySQLPowerConsumptionRepository, csvRepository domain.CSVPowerConsumptionRepository) PowerConsumptionService {
	return &PowerConsumptionServiceImpl{
		mysqlRepository,
		csvRepository,
	}
}

// ReduceInformation: reduce the information a only one record by group division
//
// Parameters:
// groupConsumptions: has the information matched between groups and information
//
// Returns:
// return reduced and one record by group division
func (s *PowerConsumptionServiceImpl) CheckingQueryParamConstrains(meterIDs string, kindPeriod string, startDate string, endDate string) (*domain.UserConsumptionQueryParams, error) {
	var numberArrayMeterIDs []int
	timeStartDate, err := domain.StrToDate(startDate)
	if err != nil {
		logrus.Errorf("Error: converting string to date startDate %s", err.Error())
		return nil, err
	}
	timeEndDate, err := domain.StrToDate(endDate)
	if err != nil {
		logrus.Errorf("Error: converting string to date endDate %s", err.Error())
		return nil, err
	}
	if timeStartDate.After(timeEndDate) {
		logrus.Errorf("Error: Invalid dates, start date must be before end date %s %s", startDate, endDate)
		return nil, fmt.Errorf("Error: Invalid dates, start date must be before end date %s %s", startDate, endDate)
	}
	timeEndDateMidnight := timeEndDate.AddDate(0, 0, 1).Add(-time.Second)

	arraymeterIDs := strings.Split(meterIDs, ",")
	if len(arraymeterIDs) == 0 {
		logrus.Errorf("Error: the array is empty %s", err.Error())
		return nil, err
	}

	for _, meterID := range arraymeterIDs {
		numberMeterID, err := domain.StrToInt(meterID)
		if err != nil {
			logrus.Errorf("Error: converting str to int meterID %s", err.Error())
			return nil, err
		}
		numberArrayMeterIDs = append(numberArrayMeterIDs, numberMeterID)
	}

	checkedKindPeriod, err := s.ChekingKindPeriod(kindPeriod)
	if err != nil {
		logrus.Errorf("Error: cheking kind period %s", err.Error())
		return nil, err
	}
	logrus.Info("the information was succefully checked all queryparms are available")
	return &domain.UserConsumptionQueryParams{
		StartDate:  timeStartDate,
		EndDate:    timeEndDateMidnight,
		MeterIDs:   numberArrayMeterIDs,
		KindPeriod: checkedKindPeriod,
	}, nil
}

// GetConsumptionByMeterIDAndWindowTime: this function check the query params for see if everithing it's ok then
// create a go routine by meter get all the information regarding this meter in a specific window time and then
// organize the information and return it
//
// Parameters:
// meterIDs: has all meterids
// startDate: has the date to start findings
// endDate: has the date to end findings
// kindPeriod: the period of time to organize the information
//
// Returns:
// return reduced and one record by group division
func (s *PowerConsumptionServiceImpl) GetConsumptionByMeterIDAndWindowTime(meterIDs, startDate, endDate, kindPeriod string) ([]Serializer, error) {

	chekedQueryParams, err := s.CheckingQueryParamConstrains(meterIDs, kindPeriod, startDate, endDate)
	if err != nil {
		return nil, err
	}
	userConsumptionChannel := make(chan Serializer, len(chekedQueryParams.MeterIDs))
	errorUserConsumptionChannel := make(chan error, len(chekedQueryParams.MeterIDs))
	wg := sync.WaitGroup{}
	if err != nil {
		return nil, err
	}

	for _, meterID := range chekedQueryParams.MeterIDs {
		wg.Add(1)
		go func(meterID int) {
			getInformation, err := s.mysqlRepository.GetConsumptionByMeterIDAndWindowTime(chekedQueryParams.StartDate, chekedQueryParams.EndDate, meterID)
			defer wg.Done()
			if err != nil {
				logrus.Errorf("Error geting the information %s meterID %d", err.Error(), meterID)
				errorUserConsumptionChannel <- err
				return
			}
			filter := NewFilter(chekedQueryParams.KindPeriod, chekedQueryParams.StartDate, chekedQueryParams.EndDate, getInformation)
			serializer := GetConsumptionData(filter)
			serializer.MeterID = meterID
			userConsumptionChannel <- serializer

		}(meterID)
	}

	go func() {
		wg.Wait()
		close(userConsumptionChannel)
		close(errorUserConsumptionChannel)
	}()

	var allUserConsumptions []Serializer
	for userConsumption := range userConsumptionChannel {
		allUserConsumptions = append(allUserConsumptions, userConsumption)
	}

	err = <-errorUserConsumptionChannel
	if err != nil {
		return nil, err
	}

	return allUserConsumptions, nil
}

// ChekingKindPeriod: this function check if the kind of period is allowed
//
// Parameters:
// meterIDs: has all meterids
// startDate: has the date to start findings
// endDate: has the date to end findings
// kindPeriod: the period of time to organize the information
//
// Returns:
// return reduced and one record by group division
func (s *PowerConsumptionServiceImpl) ChekingKindPeriod(kindPeriod string) (string, error) {
	lowerCaseKindPeriod := strings.ToLower(kindPeriod)
	trimAndLowerCaseKindPeriod := strings.Trim(lowerCaseKindPeriod, " ")
	switch trimAndLowerCaseKindPeriod {
	case constants.PeriodKindMonthly:
		return trimAndLowerCaseKindPeriod, nil
	case constants.PeriodKindWeekly:
		return trimAndLowerCaseKindPeriod, nil
	case constants.PeriodKindDaily:
		return trimAndLowerCaseKindPeriod, nil
	default:
		return "", fmt.Errorf("Error: kind period not allowed %s", trimAndLowerCaseKindPeriod)
	}
}

// ImportCsvToDatabase: this function convert and multipart file with extension csv to struct then push the information
// in the database
//
// Parameters:
// file
//
// Returns:
// return and error if the function fails or nil if it's not
func (s *PowerConsumptionServiceImpl) ImportCsvToDatabase(file *multipart.File) error {
	csvUsersConsumption, err := s.csvRepository.ConvertCSVToStruct(file)
	var usersConsumption []*domain.UserConsumption
	if err != nil {
		return err
	}

	for _, csvUserConsumption := range csvUsersConsumption {
		userConsumption, err := csvUserConsumption.ToUserConsumption()
		if err != nil {
			return err
		}
		usersConsumption = append(usersConsumption, userConsumption)
	}

	return s.mysqlRepository.CreatePowerConsumptionRecords(usersConsumption)
}
