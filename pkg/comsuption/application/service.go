package application

import (
	"mime/multipart"

	"github.com/jeffleon1/consumption-ms/pkg/comsuption/domain"
)

type PowerConsumptionService interface {
	GetConsumptionByMeterIDAndWindowTime(startDate string, endDate string) ([]domain.UserConsumption, error)
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

func (s *PowerConsumptionServiceImpl) GetConsumptionByMeterIDAndWindowTime(startDate string, endDate string) ([]domain.UserConsumption, error) {
	return s.postgresRepository.GetConsumptionByMeterIDAndWindowTime(startDate, endDate)
}

func (s *PowerConsumptionServiceImpl) ImportCsvToDatabase(file *multipart.File) error {
	usersConsumption, err := s.csvRepository.ConvertCSVToStruct(file)
	if err != nil {
		return err
	}

	return s.postgresRepository.CreatePowerConsumptionRecords(usersConsumption)
}
