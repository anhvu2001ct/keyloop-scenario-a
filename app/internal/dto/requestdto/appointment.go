package requestdto

import "time"

type BookAppointmentRequest struct {
	CustomerID    int64     `json:"customer_id" validate:"required"`
	DealershipID  int64     `json:"dealership_id" validate:"required"`
	ServiceTypeID int64     `json:"service_type_id" validate:"required"`
	TechnicianID  int64     `json:"technician_id" validate:"required"`
	VehicleID     *int64    `json:"vehicle_id"`
	StartAt       time.Time `json:"start_at" validate:"required"`
}

type CancelAppointmentRequest struct {
	Description string `json:"description"`
}

type CompleteAppointmentRequest struct {
	Description string `json:"description"`
}
