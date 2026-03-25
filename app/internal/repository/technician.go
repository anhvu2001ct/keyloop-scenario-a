package repository

import (
	"context"
	"errors"
	"fmt"
	"scenario-a/internal/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Technician interface {
	CommonRepository[*model.Technician]
	CheckAvailableForUpdate(
		ctx context.Context,
		startAt, endAt time.Time,
		technicianID, serviceTypeID int64,
	) (bool, error)
}

type technicianImpl struct {
	CommonRepository[*model.Technician]
}

func NewTechnician(gormDB *gorm.DB) Technician {
	return &technicianImpl{
		CommonRepository: NewCommonRepo[*model.Technician](gormDB),
	}
}

func (repo *technicianImpl) CheckAvailableForUpdate(
	ctx context.Context,
	startAt, endAt time.Time,
	technicianID, serviceTypeID int64,
) (bool, error) {
	db := repo.GetDB(ctx)

	var available int
	err := db.Model(&model.Technician{}).
		Select("1").
		Where("id = ?", technicianID).
		Where("EXISTS (?)", db.Model(&model.TechnicianServiceType{}).
			Select("1").
			Where("technician_id = ? AND service_type_id = ?", technicianID, serviceTypeID),
		).
		Where("NOT EXISTS (?)",
			db.Model(&model.Appointment{}).
				Select("1").
				Where("status != ?", model.AppointmentStatusCancelled).
				Where("start_at < ? AND end_at > ?", endAt, startAt).
				Where("technician_id = technicians.id")).
		Clauses(clause.Locking{Strength: "UPDATE", Table: clause.Table{Name: clause.CurrentTable}}).
		Take(&available).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	if err != nil {
		err = fmt.Errorf("failed to check technician availability: %w", err)
	}

	return available == 1, err
}
