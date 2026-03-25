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

type ServiceBay interface {
	CommonRepository[*model.ServiceBay]
	GetAvailableIDForUpdate(
		ctx context.Context,
		startAt, endAt time.Time,
		dealershipID int64,
	) (int64, error)
}

type serviceBayImpl struct {
	CommonRepository[*model.ServiceBay]
}

func NewServiceBay(gormDB *gorm.DB) ServiceBay {
	return &serviceBayImpl{
		CommonRepository: NewCommonRepo[*model.ServiceBay](gormDB),
	}
}

func (repo *serviceBayImpl) GetAvailableIDForUpdate(
	ctx context.Context,
	startAt, endAt time.Time,
	dealershipID int64,
) (int64, error) {
	db := repo.GetDB(ctx)

	var serviceBayID int64
	err := db.Model(&model.ServiceBay{}).
		Select("id").
		Where("dealership_id = ?", dealershipID).
		Where("NOT EXISTS (?)",
			db.Model(&model.Appointment{}).
				Select("1").
				Where("status != ?", model.AppointmentStatusCancelled).
				Where("start_at < ? AND end_at > ?", endAt, startAt).
				Where("service_bay_id = service_bays.id")).
		Clauses(clause.Locking{Strength: "UPDATE", Table: clause.Table{Name: clause.CurrentTable}, Options: "SKIP LOCKED"}).
		Limit(1).
		Take(&serviceBayID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	if err != nil {
		err = fmt.Errorf("failed to get available service bay: %w", err)
	}

	return serviceBayID, err
}
