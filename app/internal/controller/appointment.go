package controller

import (
	"net/http"

	"scenario-a/internal/dto/requestdto"
	"scenario-a/internal/dto/responsedto"
	"scenario-a/internal/model"
	"scenario-a/internal/repository"
	"scenario-a/internal/service"
	"scenario-a/internal/telemetry"

	"github.com/labstack/echo/v5"
	"github.com/samber/lo"
)

type Appointment struct {
	appointmentService *service.Appointment
	appointmentRepo    repository.Appointment
}

func NewAppointment(
	appointmentService *service.Appointment,
	appointmentRepo repository.Appointment,
) *Appointment {
	return &Appointment{
		appointmentService: appointmentService,
		appointmentRepo:    appointmentRepo,
	}
}

func (ctrl *Appointment) ListAppointments(c *echo.Context) error {
	ctx, span := telemetry.Tracer.Start(c.Request().Context(), "ListAppointments")
	defer span.End()

	appointments, err := ctrl.appointmentRepo.FindAll(ctx)
	if err != nil {
		return sendError(c, err)
	}

	resp := lo.Map(appointments, func(appointment *model.Appointment, _ int) *responsedto.Appointment {
		return (&responsedto.Appointment{}).FromModel(appointment)
	})

	return c.JSON(http.StatusOK, &responsedto.ListAppointmentsResponse{
		Size:  len(resp),
		Items: resp,
	})
}

func (ctrl *Appointment) BookAppointment(c *echo.Context) error {
	req := &requestdto.BookAppointmentRequest{}
	if err := bindBody(c, req); err != nil {
		return sendError(c, err)
	}

	if err := validateStruct(req); err != nil {
		return sendError(c, err)
	}

	resp, err := ctrl.appointmentService.Book(c.Request().Context(), req)
	if err != nil {
		return sendError(c, err)
	}

	return c.JSON(http.StatusOK, resp)
}

func (ctrl *Appointment) CancelAppointment(c *echo.Context) error {
	uuid, err := getPathParam[string](c, "uuid")
	if err != nil {
		return sendError(c, err)
	}

	req := &requestdto.CancelAppointmentRequest{}
	if err := bindBody(c, req); err != nil {
		return sendError(c, err)
	}

	if err := validateStruct(req); err != nil {
		return sendError(c, err)
	}

	resp, err := ctrl.appointmentService.Cancel(c.Request().Context(), uuid, req)
	if err != nil {
		return sendError(c, err)
	}

	return c.JSON(http.StatusOK, resp)
}

func (ctrl *Appointment) CompleteAppointment(c *echo.Context) error {
	uuid, err := getPathParam[string](c, "uuid")
	if err != nil {
		return sendError(c, err)
	}

	req := &requestdto.CompleteAppointmentRequest{}
	if err := bindBody(c, req); err != nil {
		return sendError(c, err)
	}

	if err := validateStruct(req); err != nil {
		return sendError(c, err)
	}

	resp, err := ctrl.appointmentService.Complete(c.Request().Context(), uuid, req)
	if err != nil {
		return sendError(c, err)
	}

	return c.JSON(http.StatusOK, resp)
}
