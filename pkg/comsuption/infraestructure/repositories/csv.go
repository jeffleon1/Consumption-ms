package repositories

import (
	"mime/multipart"

	"github.com/gocarina/gocsv"
	"github.com/jeffleon1/consumption-ms/pkg/comsuption/domain"
)

type CSVConsumptionRepositoryImpl struct{}

func NewCSVConsumptionRepository() domain.CSVPowerConsumptionRepository {
	return &CSVConsumptionRepositoryImpl{}
}

func (c *CSVConsumptionRepositoryImpl) ConvertCSVToStruct(file *multipart.File) ([]*domain.UserConsumption, error) {
	var userConsumption []*domain.UserConsumption
	if err := gocsv.UnmarshalMultipartFile(file, &userConsumption); err != nil {
		return nil, err
	}

	return userConsumption, nil
}
