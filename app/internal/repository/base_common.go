package repository

import (
	"context"
	"scenario-a/internal/model"

	"gorm.io/gorm"
)

type CommonRepository[T model.CommonModeler] interface {
	BaseRepository[T]
	FindByUUID(ctx context.Context, uuid string) (T, error)
}

type commonRepositoryImpl[T model.CommonModeler] struct {
	*baseRepoImpl[T]
}

func NewCommonRepo[T model.CommonModeler](gormDB *gorm.DB) CommonRepository[T] {
	return &commonRepositoryImpl[T]{
		baseRepoImpl: NewBaseRepo[T](gormDB).(*baseRepoImpl[T]),
	}
}

func (repo *commonRepositoryImpl[T]) FindByUUID(ctx context.Context, uuid string) (T, error) {
	return repo.FindOneBy(ctx, map[string]any{"uuid": uuid})
}
