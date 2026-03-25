package model

import "time"

type AppointmentStatus string

const (
	AppointmentStatusCreated   AppointmentStatus = "created"
	AppointmentStatusCompleted AppointmentStatus = "completed"
	AppointmentStatusCancelled AppointmentStatus = "cancelled"
)

type Appointment struct {
	CommonModel

	CustomerID    int64
	DealershipID  int64
	ServiceBayID  int64
	TechnicianID  int64
	ServiceTypeID int64
	VehicleID     *int64
	Status        AppointmentStatus
	StartAt       time.Time
	EndAt         time.Time
	Description   *string
}

func (*Appointment) TableName() string {
	return "appointments"
}

func (a *Appointment) CanCancel() bool {
	return a.Status == AppointmentStatusCreated
}

func (a *Appointment) CanComplete() bool {
	return a.Status == AppointmentStatusCreated
}
