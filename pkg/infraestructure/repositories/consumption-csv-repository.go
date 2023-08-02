package repositories

import (
	"mime/multipart"

	"github.com/gocarina/gocsv"
	"github.com/jeffleon1/consumption-ms/pkg/domain"
	"github.com/sirupsen/logrus"
)

type CSVConsumptionRepositoryImpl struct{}

func NewCSVConsumptionRepository() domain.CSVPowerConsumptionRepository {
	return &CSVConsumptionRepositoryImpl{}
}

func (c *CSVConsumptionRepositoryImpl) ConvertCSVToStruct(file *multipart.File) ([]*domain.CSVUserConsumption, error) {
	var userConsumption []*domain.CSVUserConsumption
	if err := gocsv.UnmarshalMultipartFile(file, &userConsumption); err != nil {
		logrus.Errorf("Error while converting from csv to structure %s", err.Error())
		return nil, err
	}
	logrus.Info("csv to struct conversion successfully performed")
	return userConsumption, nil
}
