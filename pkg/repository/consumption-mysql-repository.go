package repositories

import (
	"math"
	"time"

	"github.com/jeffleon1/consumption-ms/pkg/domain"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MySQLPowerConsumptionRepositoryImpl struct {
	db *gorm.DB
}

func NewMySQLPowerConsumptionRepository(db *gorm.DB) domain.MySQLPowerConsumptionRepository {
	return &MySQLPowerConsumptionRepositoryImpl{
		db,
	}
}

// GetConsumptionByMeterIDAndWindowTime: get all the records in a window time for a meter
//
// Parámeters:
// startDate - the start date to find the records.
// endDate - the end date to find the records.
// meterID - the meter id to find the records.
//
// Returns:
// return an array that represents the database domain
func (p *MySQLPowerConsumptionRepositoryImpl) GetConsumptionByMeterIDAndWindowTime(startDate, endDate time.Time, meterID int) ([]domain.UserConsumption, error) {
	var userPowerConsumption []domain.UserConsumption
	err := p.db.Where("date BETWEEN ? AND ? AND meter_id=?", startDate, endDate, meterID).Find(&userPowerConsumption).Error
	if err != nil {
		logrus.Errorf("Error: converting csv to struct %s", err.Error())
		return nil, err
	}
	logrus.Info("the convertion from csv to struct was succesfully")
	return userPowerConsumption, nil

}

// CreatePowerConsumptionRecords: create a records for user power consumption
//
// Parámeters:
// usersPowerConsumption - user power consumption domain.
//
// Returns:
// return an error if something goes wrong in the insertion of nil if it's not
func (p *MySQLPowerConsumptionRepositoryImpl) CreatePowerConsumptionRecords(usersPowerConsumption []*domain.UserConsumption) error {
	recordSize := len(usersPowerConsumption)
	recordLimit := 4000
	lotsNumber := int(math.Ceil(float64(recordSize) / float64(recordLimit)))

	for i := 0; i < lotsNumber; i++ {
		begin := i * recordLimit
		end := int(math.Min(float64((i+1)*recordLimit), float64(recordSize)))
		lot := usersPowerConsumption[begin:end]
		logrus.Info("Lot ", begin, end)
		errors := p.db.Create(&lot).Error
		if errors != nil {
			logrus.Errorf("Error inserting in the lot: %s", errors.Error())
			return errors
		}

		time.Sleep(500 * time.Millisecond)
	}
	logrus.Info("the Insertion was succesfully in user_consumption database")
	return nil

}

// ModelMigration: do the model migration to gorm
//
// Returns:
// return an error if something goes wrong in the migration of nil if it's not
func (p *MySQLPowerConsumptionRepositoryImpl) ModelMigration() error {
	return p.db.AutoMigrate(&domain.UserConsumption{})
}
