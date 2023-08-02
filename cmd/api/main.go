package main

import (
	"fmt"
	"os"

	config "github.com/jeffleon1/consumption-ms/internal/configuration"
	"github.com/jeffleon1/consumption-ms/pkg/application"
	"github.com/jeffleon1/consumption-ms/pkg/infraestructure"
	"github.com/jeffleon1/consumption-ms/pkg/infraestructure/repositories"
	"github.com/sirupsen/logrus"
)

// @title           Consumption API
// @version         1.0
// @description     It's and API that provides information regarding the energy consumption with three diferent filters or rates of time weekly. monthly and daily in a given window time.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @BasePath  /api/v1
// @query.collection, format multi
func main() {
	port := config.Config.APP.PORT
	db, err := config.Config.DatabaseInit()
	if err != nil {
		logrus.Fatalf("Fatal Error: the database could not connect %s", err.Error())
		os.Exit(1)
	}
	powerConsumptionMySQLRepository := repositories.NewMySQLPowerConsumptionRepository(db)
	err = powerConsumptionMySQLRepository.ModelMigration()
	if err != nil {
		logrus.Fatalf("Fatal Error: It was not possible to migrate the model %s", err.Error())
		os.Exit(1)
	}
	powerConsumptionCSVRepository := repositories.NewCSVConsumptionRepository()
	powerConsumptionService := application.NewPowerConsumptionService(powerConsumptionMySQLRepository, powerConsumptionCSVRepository)
	powerConsumptionHandler := infraestructure.NewPowerConsumptionHandler(powerConsumptionService)
	powerConsumptionRoutes := infraestructure.NewRoutes(powerConsumptionHandler)

	r := infraestructure.NewRouter(infraestructure.RoutesGroup{
		PowerConsumption: powerConsumptionRoutes,
		Swagger:          infraestructure.NewSwaggerDocsRoutes(),
	})

	logrus.Fatal(r.Run(fmt.Sprintf(":%s", port)))
}
