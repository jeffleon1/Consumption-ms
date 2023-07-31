package repositories

import (
	"math"
	"time"

	"github.com/jeffleon1/consumption-ms/pkg/comsuption/domain"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PostgreSQLPowerConsumptionRepositoryImpl struct {
	db *gorm.DB
}

func NewPostgreSQLPowerConsumptionRepository(db *gorm.DB) domain.PostgresPowerConsumptionRepository {
	return &PostgreSQLPowerConsumptionRepositoryImpl{
		db,
	}
}

func (p *PostgreSQLPowerConsumptionRepositoryImpl) GetConsumptionByMeterIDAndWindowTime(startDate string, endDate string) ([]domain.UserConsumption, error) {
	var userPowerConsumption []domain.UserConsumption
	err := p.db.Where("date BETWEEN ? AND ?", startDate, endDate).Find(&userPowerConsumption).Error
	if err != nil {
		return nil, err
	}

	return userPowerConsumption, nil

}

func (p *PostgreSQLPowerConsumptionRepositoryImpl) CreatePowerConsumptionRecords(usersPowerConsumption []*domain.UserConsumption) error {
	recordSize := len(usersPowerConsumption)
	recordLimit := 4000
	lotsNumber := int(math.Ceil(float64(recordSize) / float64(recordLimit)))

	for i := 0; i < lotsNumber; i++ {
		begin := i * recordLimit
		end := int(math.Min(float64((i+1)*recordLimit), float64(recordSize)))
		lot := usersPowerConsumption[begin:end]
		logrus.Println("Lot ", begin, end)
		errors := p.db.Create(&lot).Error
		if errors != nil {
			logrus.Errorf("Error inserting in the lot: %s", errors.Error())
			return errors
		}

		time.Sleep(500 * time.Millisecond)
	}

	return nil

}

func (p *PostgreSQLPowerConsumptionRepositoryImpl) ModelMigration() error {
	return p.db.AutoMigrate(&domain.UserConsumption{})
}
