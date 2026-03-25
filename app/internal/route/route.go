package route

import (
	"scenario-a/internal/dep"

	"github.com/labstack/echo/v5"
)

func Load(e *echo.Echo, deps *dep.Dependencies) {
	initAppointmentRoutes(e, deps)
}
