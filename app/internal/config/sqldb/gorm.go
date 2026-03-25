package sqldb

import (
	"database/sql"
	"fmt"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initGorm(db *sql.DB, enableDebug bool) *gorm.DB {
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("error while init gorm from sql db: %v", err))
	}

	if err := gormDB.Use(otelgorm.NewPlugin()); err != nil {
		panic(err)
	}

	if enableDebug {
		gormDB = gormDB.Debug()
	}

	return gormDB
}

func initGormFromDSN(dsn string, enableDebug bool) *gorm.DB {
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("error while init gorm from dsn: %v", err))
	}

	if err := gormDB.Use(otelgorm.NewPlugin()); err != nil {
		panic(err)
	}

	if enableDebug {
		gormDB = gormDB.Debug()
	}

	return gormDB
}
