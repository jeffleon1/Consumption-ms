package domain

import (
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	constants "github.com/jeffleon1/consumption-ms/internal/constans"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserConsumption struct {
	gorm.Model
	ID                 string    `gorm:"primary_key;auto_increment" json:"id" csv:"id"`
	MeterID            int       `gorm:"meter_id" json:"meter_id" csv:"meter_id"`
	ActiveEnergy       float64   `gorm:"active_energy" json:"active_energy" csv:"active_energy"`
	ReactiveEnergy     float64   `gorm:"reactive_energy" json:"reactive_energy" csv:"reactive_energy"`
	CapacitiveReactive float64   `gorm:"capacity_energy" json:"capacitive_reactive" csv:"capacitive_reactive"`
	Solar              float64   `gorm:"solar" json:"solar" csv:"solar"`
	Date               time.Time `gorm:"date" json:"date" csv:"date"`
}

type UserConsumptionQueryParams struct {
	StartDate  time.Time
	EndDate    time.Time
	MeterIDs   []int
	KindPeriod string
}

type CSVUserConsumption struct {
	ID                 string  `json:"id" csv:"id"`
	MeterID            string  `json:"meter_id" csv:"meter_id"`
	ActiveEnergy       float64 `json:"active_energy" csv:"active_energy"`
	ReactiveEnergy     float64 `json:"reactive_energy" csv:"reactive_energy"`
	CapacitiveReactive float64 `json:"capacitive_reactive" csv:"capacitive_reactive"`
	Solar              float64 `json:"solar" csv:"solar"`
	Date               string  `json:"date" csv:"date"`
}

func (u CSVUserConsumption) ToUserConsumption() (*UserConsumption, error) {
	objectDate, err := StrToDate(u.Date)
	if err != nil {
		logrus.Errorf("Error trying to cast date string to Time.time %s", u.Date)
		return nil, err
	}

	numberMeterID, err := StrToInt(u.MeterID)
	if err != nil {
		logrus.Errorf("Error trying to cast string to int %s", u.Date)
		return nil, err
	}
	return &UserConsumption{
		ID:                 u.ID,
		MeterID:            numberMeterID,
		ActiveEnergy:       u.ActiveEnergy,
		ReactiveEnergy:     u.ReactiveEnergy,
		CapacitiveReactive: u.ReactiveEnergy,
		Solar:              u.Solar,
		Date:               objectDate,
	}, nil
}

func StrToInt(meterID string) (int, error) {
	return strconv.Atoi(meterID)
}

func TimeTostr(date time.Time, format string) string {
	return date.Format(format)
}

func StrToDate(date string) (time.Time, error) {
	dateFormat := constants.DateFormatDateTimeWithTZ
	if !hasHourAndMinutes(date) {
		date += " 00:00:00+00"
	}
	dateFormated, err := time.Parse(dateFormat, date)
	if err != nil {
		logrus.Errorf("Error Parsing date %s", err.Error())
		return time.Time{}, err
	}
	logrus.Info("Parsing Successfully do it ")
	return dateFormated, nil

}

func hasHourAndMinutes(date string) bool {
	arrayDate := strings.Split(date, " ")
	if len(arrayDate) <= 1 {
		return false
	}

	count := strings.Count(arrayDate[1], ":")
	return count == 2
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . MySQLPowerConsumptionRepository
type MySQLPowerConsumptionRepository interface {
	GetConsumptionByMeterIDAndWindowTime(startDate, endDate time.Time, meterID int) ([]UserConsumption, error)
	CreatePowerConsumptionRecords(usersPowerConsumption []*UserConsumption) error
	ModelMigration() error
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . CSVPowerConsumptionRepository
type CSVPowerConsumptionRepository interface {
	ConvertCSVToStruct(file *multipart.File) ([]*CSVUserConsumption, error)
}
