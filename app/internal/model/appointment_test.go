package model_test

import (
	"scenario-a/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppointment_CanCancel(t *testing.T) {
	tests := []struct {
		name     string
		status   model.AppointmentStatus
		expected bool
	}{
		{"created status can cancel", model.AppointmentStatusCreated, true},
		{"completed status cannot cancel", model.AppointmentStatusCompleted, false},
		{"cancelled status cannot cancel", model.AppointmentStatusCancelled, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appt := &model.Appointment{Status: tt.status}
			assert.Equal(t, tt.expected, appt.CanCancel())
		})
	}
}

func TestAppointment_CanComplete(t *testing.T) {
	tests := []struct {
		name     string
		status   model.AppointmentStatus
		expected bool
	}{
		{"created status can complete", model.AppointmentStatusCreated, true},
		{"completed status cannot complete", model.AppointmentStatusCompleted, false},
		{"cancelled status cannot complete", model.AppointmentStatusCancelled, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appt := &model.Appointment{Status: tt.status}
			assert.Equal(t, tt.expected, appt.CanComplete())
		})
	}
}

func TestAppointment_TableName(t *testing.T) {
	appt := &model.Appointment{}
	assert.Equal(t, "appointments", appt.TableName())
}
