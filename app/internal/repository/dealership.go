package repository

import (
	"scenario-a/internal/model"

	"gorm.io/gorm"
)

type Dealership interface {
	CommonRepository[*model.Dealership]
}

type dealershipImpl struct {
	CommonRepository[*model.Dealership]
}

func NewDealership(gormDB *gorm.DB) Dealership {
	return &dealershipImpl{
		CommonRepository: NewCommonRepo[*model.Dealership](gormDB),
	}
}
