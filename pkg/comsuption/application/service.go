package application

import (
	"fmt"
	"mime/multipart"
	"strings"
	"sync"

	"github.com/jeffleon1/consumption-ms/pkg/comsuption/domain"
	"github.com/sirupsen/logrus"
)

type PowerConsumptionService interface {
	GetConsumptionByMeterIDAndWindowTime(meterIDs string, kindPeriod string, startDate string, endDate string) ([]domain.UserConsumption, error)
	ImportCsvToDatabase(file *multipart.File) error
}

type PowerConsumptionServiceImpl struct {
	postgresRepository domain.PostgresPowerConsumptionRepository
	csvRepository      domain.CSVPowerConsumptionRepository
}

func NewPowerConsumptionService(postgresRepository domain.PostgresPowerConsumptionRepository, csvRepository domain.CSVPowerConsumptionRepository) PowerConsumptionService {
	return &PowerConsumptionServiceImpl{
		postgresRepository,
		csvRepository,
	}
}

func (s *PowerConsumptionServiceImpl) GetConsumptionByMeterIDAndWindowTime(meterIDs, startDate, endDate, kindPeriod string) ([]domain.UserConsumption, error) {

	chekedQueryParams, err := s.checkingQueryParamConstrains(meterIDs, kindPeriod, startDate, endDate)
	userConsumptionChannel := make(chan []domain.UserConsumption, len(chekedQueryParams.MeterIDs))
	wg := sync.WaitGroup{}
	if err != nil {
		return nil, err
	}

	wg.Add(len(meterIDs))
	for _, meterID := range chekedQueryParams.MeterIDs {
		go func(meterID int) {
			getInformation, _ := s.postgresRepository.GetConsumptionByMeterIDAndWindowTime(chekedQueryParams.StartDate, chekedQueryParams.EndDate, meterID)
			userConsumptionChannel <- getInformation
			defer wg.Done()
			fmt.Println(getInformation)
		}(meterID)
	}
	wg.Wait()
	return nil, nil
}

func (s *PowerConsumptionServiceImpl) checkingQueryParamConstrains(meterIDs string, kindPeriod string, startDate string, endDate string) (*domain.UserConsumptionQueryParams, error) {
	timeStartDate, err := domain.StrToDate(startDate)
	var numberArrayMeterIDs []int
	if err != nil {
		logrus.Errorf("Error: %s", err.Error())
		return nil, err
	}
	timeEndDate, err := domain.StrToDate(endDate)
	if err != nil {
		logrus.Errorf("Error: %s", err.Error())
		return nil, err
	}

	arraymeterIDs := strings.Split(meterIDs, ",")
	if len(arraymeterIDs) == 0 {
		logrus.Errorf("Error: %s", err.Error())
		return nil, err
	}

	for _, meterID := range arraymeterIDs {
		numberMeterID, err := domain.StrToInt(meterID)
		if err != nil {
			logrus.Errorf("Error: %s", err.Error())
			return nil, err
		}
		numberArrayMeterIDs = append(numberArrayMeterIDs, numberMeterID)
	}

	checkedKindPeriod, err := s.chekingKindPeriod(kindPeriod)
	if err != nil {
		logrus.Errorf("Error: %s", err.Error())
		return nil, err
	}

	return &domain.UserConsumptionQueryParams{
		StartDate:  timeStartDate,
		EndDate:    timeEndDate,
		MeterIDs:   numberArrayMeterIDs,
		KindPeriod: checkedKindPeriod,
	}, nil
}

func (s *PowerConsumptionServiceImpl) chekingKindPeriod(kindPeriod string) (string, error) {
	lowerCaseKindPeriod := strings.ToLower(kindPeriod)
	trimAndLowerCaseKindPeriod := strings.Trim(lowerCaseKindPeriod, " ")
	switch trimAndLowerCaseKindPeriod {
	case "monthly":
		return trimAndLowerCaseKindPeriod, nil
	case "weekly":
		return trimAndLowerCaseKindPeriod, nil
	case "daily":
		return trimAndLowerCaseKindPeriod, nil
	default:
		return "", fmt.Errorf("Error kind period not allowed %s", trimAndLowerCaseKindPeriod)
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

	return s.postgresRepository.CreatePowerConsumptionRecords(usersConsumption)
}
