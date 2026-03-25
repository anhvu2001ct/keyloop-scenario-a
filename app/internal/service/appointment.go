package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"scenario-a/internal/dto/requestdto"
	"scenario-a/internal/dto/responsedto"
	"scenario-a/internal/model"
	"scenario-a/internal/repository"
	"scenario-a/internal/telemetry"
	"scenario-a/pkg/common"
)

type Appointment struct {
	appointmentRepo repository.Appointment
	serviceTypeRepo repository.ServiceType
	dealershipRepo  repository.Dealership
	technicianRepo  repository.Technician
	serviceBayRepo  repository.ServiceBay
	customerRepo    repository.Customer
}

func NewAppointment(
	appointmentRepo repository.Appointment,
	serviceTypeRepo repository.ServiceType,
	dealershipRepo repository.Dealership,
	technicianRepo repository.Technician,
	serviceBayRepo repository.ServiceBay,
	customerRepo repository.Customer,
) *Appointment {
	return &Appointment{
		appointmentRepo: appointmentRepo,
		serviceTypeRepo: serviceTypeRepo,
		dealershipRepo:  dealershipRepo,
		technicianRepo:  technicianRepo,
		serviceBayRepo:  serviceBayRepo,
		customerRepo:    customerRepo,
	}
}

func (s *Appointment) Book(
	ctx context.Context,
	req *requestdto.BookAppointmentRequest,
) (*responsedto.Appointment, error) {
	serviceType, err := s.serviceTypeRepo.FindByID(ctx, req.ServiceTypeID)
	if err != nil {
		return nil, err
	}

	_, err = s.customerRepo.FindByID(ctx, req.CustomerID)
	if err != nil {
		return nil, err
	}

	dealership, err := s.dealershipRepo.FindByID(ctx, req.DealershipID)
	if err != nil {
		return nil, err
	}

	technician, err := s.technicianRepo.FindByID(ctx, req.TechnicianID)
	if err != nil {
		return nil, err
	}

	if technician.DealershipID != dealership.ID {
		return nil, &common.AppError{
			Message: "The requested technician does not belong to the requested dealership",
			Type:    common.AppErrorValidationFailed,
		}
	}

	endAt := req.StartAt.Add(time.Duration(serviceType.DurationMinutes) * time.Minute)

	if err := dealership.IsOpen(req.StartAt, endAt); err != nil {
		return nil, &common.AppError{
			Cause: err,
			Type:  common.AppErrorValidationFailed,
		}
	}

	var newAppointment *model.Appointment

	ctx, span := telemetry.Tracer.Start(ctx, "CheckAvailabilityConflict")
	defer span.End()

	err = s.appointmentRepo.Transaction(ctx, func(ctx context.Context, _ *gorm.DB) error {
		isAvailable, err := s.technicianRepo.CheckAvailableForUpdate(ctx, req.StartAt, endAt, req.TechnicianID, req.ServiceTypeID)
		if err != nil {
			return err
		}
		if !isAvailable {
			return &common.AppError{
				Message: "The requested technician is not available or doesn't have the required skills",
				Type:    common.AppErrorValidationFailed,
			}
		}

		serviceBayID, err := s.serviceBayRepo.GetAvailableIDForUpdate(ctx, req.StartAt, endAt, req.DealershipID)
		if err != nil {
			return err
		}
		if serviceBayID == 0 {
			return &common.AppError{
				Message: "The requested dealership does not have any available service bays in this moment",
				Type:    common.AppErrorValidationFailed,
			}
		}

		newAppointment = &model.Appointment{
			CommonModel: model.CommonModel{
				UUID: uuid.New().String(),
			},
			CustomerID:    req.CustomerID,
			DealershipID:  req.DealershipID,
			ServiceBayID:  serviceBayID,
			TechnicianID:  req.TechnicianID,
			ServiceTypeID: req.ServiceTypeID,
			VehicleID:     req.VehicleID,
			Status:        model.AppointmentStatusCreated,
			StartAt:       req.StartAt,
			EndAt:         endAt,
		}

		if err := s.appointmentRepo.Create(ctx, newAppointment); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return (&responsedto.Appointment{}).FromModel(newAppointment), nil
}

func (s *Appointment) Cancel(
	ctx context.Context,
	appointmentUUID string,
	req *requestdto.CancelAppointmentRequest,
) (*responsedto.Appointment, error) {
	appointment, err := s.appointmentRepo.FindByUUID(ctx, appointmentUUID)
	if err != nil {
		return nil, err
	}

	if !appointment.CanCancel() {
		return nil, &common.AppError{
			Message: "The requested appointment is not in a cancellable state",
			Type:    common.AppErrorValidationFailed,
		}
	}

	if err := s.appointmentRepo.UpdateModelByMap(ctx, appointment, map[string]any{
		"status":      model.AppointmentStatusCancelled,
		"description": req.Description,
	}); err != nil {
		return nil, err
	}

	return (&responsedto.Appointment{}).FromModel(appointment), nil
}

func (s *Appointment) Complete(
	ctx context.Context,
	appointmentUUID string,
	req *requestdto.CompleteAppointmentRequest,
) (*responsedto.Appointment, error) {
	appointment, err := s.appointmentRepo.FindByUUID(ctx, appointmentUUID)
	if err != nil {
		return nil, err
	}

	if !appointment.CanComplete() {
		return nil, &common.AppError{
			Message: "The requested appointment is not in a completable state",
			Type:    common.AppErrorValidationFailed,
		}
	}

	if err := s.appointmentRepo.UpdateModelByMap(ctx, appointment, map[string]any{
		"status":      model.AppointmentStatusCompleted,
		"description": req.Description,
	}); err != nil {
		return nil, err
	}

	return (&responsedto.Appointment{}).FromModel(appointment), nil
}
