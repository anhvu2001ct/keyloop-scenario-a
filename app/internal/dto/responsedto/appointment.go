package responsedto

import (
	"scenario-a/internal/model"
	"time"
)

type Appointment struct {
	UUID        string  `json:"uuid"`
	Status      string  `json:"status"`
	Description *string `json:"description"`
	StartAt     string  `json:"start_at"`
	EndAt       string  `json:"end_at"`
}

func (r *Appointment) FromModel(appointment *model.Appointment) *Appointment {
	r.UUID = appointment.UUID
	r.Status = string(appointment.Status)
	r.Description = appointment.Description
	r.StartAt = appointment.StartAt.Format(time.RFC3339)
	r.EndAt = appointment.EndAt.Format(time.RFC3339)
	return r
}

type ListAppointmentsResponse struct {
	Size  int            `json:"size"`
	Items []*Appointment `json:"items"`
}
