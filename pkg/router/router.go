package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jeffleon1/consumption-ms/pkg/comsuption/infraestructure"
)

type Router interface {
	Run(addr ...string) error
}

func NewRouter(routes RoutesGroup) Router {
	route := gin.Default()
	public := route.Group("/api/v1")
	routes.PowerConsumption.RegisterRoutes(public)
	return route
}

type RoutesGroup struct {
	PowerConsumption *infraestructure.PowerConsumptionRoutes
}
