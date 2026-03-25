package route

import (
	"scenario-a/internal/controller"
	"scenario-a/internal/dep"

	"github.com/labstack/echo/v5"
)

func initAppointmentRoutes(e *echo.Echo, deps *dep.Dependencies) {
	controller := controller.NewAppointment(deps.Service.Appointment, deps.Repository.Appointment)
	g := e.Group("/appointments")

	g.GET("", controller.ListAppointments)
	g.POST("", controller.BookAppointment)
	g.POST("/:uuid/cancel", controller.CancelAppointment)
	g.POST("/:uuid/complete", controller.CompleteAppointment)
}
