package repository

import (
	"context"

	"gorm.io/gorm"
)

type gormTxContext struct{}

var gormTxContextKey = gormTxContext{}

func GormTxFromContext(ctx context.Context) (*gorm.DB, bool) {
	val, ok := ctx.Value(gormTxContextKey).(*gorm.DB)
	return val, ok
}

func GormTxToContext(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, gormTxContextKey, db)
}
