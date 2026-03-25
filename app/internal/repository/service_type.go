package repository

import (
	"scenario-a/internal/model"

	"gorm.io/gorm"
)

type ServiceType interface {
	CommonRepository[*model.ServiceType]
}

type serviceTypeImpl struct {
	CommonRepository[*model.ServiceType]
}

func NewServiceType(gormDB *gorm.DB) ServiceType {
	return &serviceTypeImpl{
		CommonRepository: NewCommonRepo[*model.ServiceType](gormDB),
	}
}
