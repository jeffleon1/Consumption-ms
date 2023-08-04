package repositories

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jeffleon1/consumption-ms/pkg/domain"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _ = Describe("GetConsumptionByMeterIDAndWindowTime", func() {
	var (
		mockDB         *gorm.DB
		mock           sqlmock.Sqlmock
		mockDb         *sql.DB
		repositoryImpl *MySQLPowerConsumptionRepositoryImpl
		err            error
	)

	BeforeEach(func() {
		mockDb, mock, _ = sqlmock.New()
		mockDB, err = gorm.Open(mysql.New(mysql.Config{
			Conn:                      mockDb,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		repositoryImpl = &MySQLPowerConsumptionRepositoryImpl{
			db: mockDB,
		}
	})

	Context("when data is not found in the database", func() {
		It("should return an empty slice", func() {

			now := time.Now()
			startDate := now.AddDate(0, -6, 0)
			endDate := startDate.AddDate(0, 0, 7)
			meterID := 1

			columns := []string{"id", "meter_id", "active_energy", "reactive_energy", "capacitive_reactive", "solar", "date"}
			rows := sqlmock.NewRows(columns).AddRow("993c644a-bace-4e1e-9b42-2a9c18aadcfe", 1, 17234.731809999997, 10516.07749, 1, 0.6388935728544158, startDate)

			mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

			result, err := repositoryImpl.GetConsumptionByMeterIDAndWindowTime(startDate, endDate, meterID)
			Expect(err).To(BeNil())
			Expect(result).To(Not(BeNil()))

		})
	})

	Context("when data is found in the database", func() {
		It("should return the consumption data", func() {

			now := time.Now()
			startDate := now.AddDate(0, -6, 0)
			endDate := startDate.AddDate(0, 0, 7)
			meterID := 1

			mock.ExpectQuery(`SELECT`).WillReturnError(sqlmock.ErrCancelled)

			result, err := repositoryImpl.GetConsumptionByMeterIDAndWindowTime(startDate, endDate, meterID)
			Expect(err).To(Not(BeNil()))
			Expect(result).To(BeNil())

		})
	})

	Context("when data is found in the database", func() {
		It("should return the consumption data", func() {

			now := time.Now()
			startDate := now.AddDate(0, -6, 0)
			endDate := startDate.AddDate(0, 0, 7)
			meterID := 1

			rows := sqlmock.NewRows([]string{"id", "meter_id", "active_energy", "reactive_energy", "capacitive_reactive", "solar", "date"})

			mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

			result, err := repositoryImpl.GetConsumptionByMeterIDAndWindowTime(startDate, endDate, meterID)
			Expect(err).To(BeNil())
			Expect(result).To(HaveLen(0))

		})
	})
})

var _ = Describe("CreatePowerConsumptionRecords", func() {
	var (
		mockDB           *gorm.DB
		mock             sqlmock.Sqlmock
		mockDb           *sql.DB
		repositoryImpl   *MySQLPowerConsumptionRepositoryImpl
		err              error
		userConsumptions = []*domain.UserConsumption{
			{
				ID:                 "2131232131231",
				MeterID:            1,
				ActiveEnergy:       100.0,
				ReactiveEnergy:     50.0,
				CapacitiveReactive: 20.0,
				Solar:              30.0,
				Date:               time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			},
		}
	)

	BeforeEach(func() {
		mockDb, mock, _ = sqlmock.New()
		mockDB, err = gorm.Open(mysql.New(mysql.Config{
			Conn:                      mockDb,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		repositoryImpl = &MySQLPowerConsumptionRepositoryImpl{
			db: mockDB,
		}
	})

	Context("when data is create in the database", func() {
		It("should return an empty error", func() {
			num, _ := strconv.Atoi(userConsumptions[0].ID)
			mock.ExpectBegin()
			mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(int64(num), 1))
			mock.ExpectCommit()
			err := repositoryImpl.CreatePowerConsumptionRecords(userConsumptions)
			Expect(err).To(BeNil())
		})
	})

	Context("when an error occurs during insertion", func() {
		It("should return an error", func() {
			num, _ := strconv.Atoi(userConsumptions[0].ID)
			mock.ExpectBegin()
			mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(int64(num), 1))
			mock.ExpectRollback()

			err := repositoryImpl.CreatePowerConsumptionRecords(userConsumptions)

			Expect(err).ToNot(BeNil())
		})
	})
})
