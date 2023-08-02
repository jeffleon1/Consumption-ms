package application

import (
	"fmt"
	"mime/multipart"
	"strings"
	"sync"
	"time"

	constants "github.com/jeffleon1/consumption-ms/internal/const"
	"github.com/jeffleon1/consumption-ms/pkg/comsuption/domain"
	"github.com/sirupsen/logrus"
)

type PowerConsumptionService interface {
	GetConsumptionByMeterIDAndWindowTime(meterIDs string, kindPeriod string, startDate string, endDate string) ([]Serializer, error)
	ImportCsvToDatabase(file *multipart.File) error
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

func (s *PowerConsumptionServiceImpl) GetConsumptionByMeterIDAndWindowTime(meterIDs, startDate, endDate, kindPeriod string) ([]Serializer, error) {

	chekedQueryParams, err := s.checkingQueryParamConstrains(meterIDs, kindPeriod, startDate, endDate)
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

func (s *PowerConsumptionServiceImpl) checkingQueryParamConstrains(meterIDs string, kindPeriod string, startDate string, endDate string) (*domain.UserConsumptionQueryParams, error) {
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

	checkedKindPeriod, err := s.chekingKindPeriod(kindPeriod)
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

func (s *PowerConsumptionServiceImpl) chekingKindPeriod(kindPeriod string) (string, error) {
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
