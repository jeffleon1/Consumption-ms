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
	port := config.Config.APP
	powerConsumptionPostgresRepository := repositories.NewPostgreSQLPowerConsumptionRepository()
	powerConsumptionService := application.NewPowerConsumptionService(&powerConsumptionPostgresRepository)
	powerConsumptionHandler := infraestructure.NewPowerConsumptionHandler(&powerConsumptionService)
	powerConsumptionRoutes := infraestructure.NewRoutes(&powerConsumptionHandler)

	r := router.NewRouter(router.RoutesGroup{
		PowerConsumption: powerConsumptionRoutes,
	})

	logrus.Fatal(r.Run(fmt.Sprintf(":%s", port)))
}
