package repository

import (
	"scenario-a/internal/model"

	"gorm.io/gorm"
)

type Appointment interface {
	CommonRepository[*model.Appointment]
}

type appointmentImpl struct {
	CommonRepository[*model.Appointment]
}

func NewAppointment(gormDB *gorm.DB) Appointment {
	return &appointmentImpl{
		CommonRepository: NewCommonRepo[*model.Appointment](gormDB),
	}
}
