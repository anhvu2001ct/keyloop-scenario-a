package model

import "time"

type BaseModeler interface {
	GetID() int64
	TableName() string
}

type CommonModeler interface {
	BaseModeler
	GetUUID() string
}

type SoftDeletable interface {
	GetDeletedAt() *time.Time
}

type BaseModel struct {
	ID int64

	// use db default auto create/update time
	CreatedAt time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
}

func (b *BaseModel) GetID() int64 {
	return b.ID
}

type CommonModel struct {
	BaseModel

	UUID string
}

func (c *CommonModel) GetUUID() string {
	return c.UUID
}

type WithSoftDelete struct {
	DeletedAt *time.Time
}

func (w *WithSoftDelete) GetDeletedAt() *time.Time {
	return w.DeletedAt
}
