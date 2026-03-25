package repository

import (
	"scenario-a/internal/model"

	"gorm.io/gorm"
)

type Customer interface {
	CommonRepository[*model.Customer]
}

type customerImpl struct {
	CommonRepository[*model.Customer]
}

func NewCustomer(gormDB *gorm.DB) Customer {
	return &customerImpl{
		CommonRepository: NewCommonRepo[*model.Customer](gormDB),
	}
}
