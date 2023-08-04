package infraestructure

import (
	"github.com/gin-gonic/gin"
)

type Router interface {
	Run(addr ...string) error
}

func NewRouter(routes RoutesGroup) Router {
	route := gin.Default()
	public := route.Group("/api/v1")
	routes.Swagger.RegisterRoutes(public)
	routes.PowerConsumption.RegisterRoutes(public)
	return route
}

type RoutesGroup struct {
	PowerConsumption *PowerConsumptionRoutes
	Swagger          *SwaggerRoutes
}
