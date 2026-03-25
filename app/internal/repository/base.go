package repository

import (
	"context"
	"fmt"
	"scenario-a/internal/errors/dberr"
	"scenario-a/internal/model"

	"gorm.io/gorm"
)

type BaseRepository[T model.BaseModeler] interface {
	GetDB(ctx context.Context) *gorm.DB
	Transaction(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) error
	FindAll(ctx context.Context) ([]T, error)
	FindOneBy(ctx context.Context, conditions map[string]any) (T, error)
	FindByID(ctx context.Context, id int64) (T, error)
	Create(ctx context.Context, record T) error
	UpdateModelByMap(ctx context.Context, model T, updates map[string]any) error
}

type baseRepoImpl[T model.BaseModeler] struct {
	db *gorm.DB
}

func NewBaseRepo[T model.BaseModeler](gormDB *gorm.DB) BaseRepository[T] {
	return &baseRepoImpl[T]{
		db: gormDB,
	}
}

func (repo *baseRepoImpl[T]) GetDB(ctx context.Context) *gorm.DB {
	if db, ok := GormTxFromContext(ctx); ok {
		return db
	}
	return repo.db.WithContext(ctx)
}

func (repo *baseRepoImpl[T]) Transaction(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) error {
	return repo.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(GormTxToContext(ctx, tx), tx)
	})
}

func (repo *baseRepoImpl[T]) FindAll(ctx context.Context) ([]T, error) {
	return gorm.G[T](repo.GetDB(ctx)).Find(ctx)
}

func (repo *baseRepoImpl[T]) FindOneBy(ctx context.Context, conditions map[string]any) (T, error) {
	record, err := gorm.G[T](repo.GetDB(ctx)).Where(conditions).Take(ctx)
	if err == gorm.ErrRecordNotFound {
		err = &dberr.RecordNotFound{
			TableName:  record.TableName(),
			Conditions: conditions,
		}
	}
	return record, err
}

func (repo *baseRepoImpl[T]) FindByID(ctx context.Context, id int64) (T, error) {
	return repo.FindOneBy(ctx, map[string]any{"id": id})
}

func (repo *baseRepoImpl[T]) Create(ctx context.Context, record T) error {
	err := gorm.G[T](repo.GetDB(ctx)).Create(ctx, &record)
	if err != nil {
		err = fmt.Errorf("failed to create %s: %w", record.TableName(), err)
	}
	return err
}

func (repo *baseRepoImpl[T]) UpdateModelByMap(ctx context.Context, model T, updates map[string]any) error {
	return repo.GetDB(ctx).Model(model).Updates(updates).Error
}
