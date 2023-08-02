package main

import (
	"fmt"

	"github.com/jeffleon1/consumption-ms/internal/config"
	"github.com/jeffleon1/consumption-ms/pkg/comsuption/application"
	"github.com/jeffleon1/consumption-ms/pkg/comsuption/infraestructure"
	"github.com/jeffleon1/consumption-ms/pkg/comsuption/infraestructure/repositories"
	"github.com/jeffleon1/consumption-ms/pkg/router"
	"github.com/sirupsen/logrus"
)

func main() {
	port := config.Config.APP.PORT
	db, err := config.Config.DatabaseInit()
	if err != nil {
		logrus.Fatalf("Fatal Error: the database could not connect %s", err.Error())
	}
	powerConsumptionMySQLRepository := repositories.NewMySQLPowerConsumptionRepository(db)
	err = powerConsumptionMySQLRepository.ModelMigration()
	if err != nil {
		logrus.Fatalf("Fatal Error: It was not possible to migrate the model %s", err.Error())
	}
	powerConsumptionCSVRepository := repositories.NewCSVConsumptionRepository()
	powerConsumptionService := application.NewPowerConsumptionService(powerConsumptionMySQLRepository, powerConsumptionCSVRepository)
	powerConsumptionHandler := infraestructure.NewPowerConsumptionHandler(powerConsumptionService)
	powerConsumptionRoutes := infraestructure.NewRoutes(powerConsumptionHandler)

	r := router.NewRouter(router.RoutesGroup{
		PowerConsumption: powerConsumptionRoutes,
	})

	logrus.Fatal(r.Run(fmt.Sprintf(":%s", port)))
}
