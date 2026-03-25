package sqldb

import (
	"database/sql"
	"fmt"
	"scenario-a/internal/config/env"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/gorm"
)

func getDSN(cfgEnv *env.Env) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		cfgEnv.DBHost,
		cfgEnv.DBPort,
		cfgEnv.DBUser,
		cfgEnv.DBPassword,
		cfgEnv.DBName,
		cfgEnv.DBSSLMode,
	)
}

// MustInit opens a GORM *gorm.DB and the underlying standard *sql.DB using
// the connection settings from cfg. It panics on any error.
func MustInit(cfgEnv *env.Env) (*sql.DB, *gorm.DB) {
	dsn := getDSN(cfgEnv)
	gormDB := initGormFromDSN(dsn, cfgEnv.AppEnv != env.Production)

	sqlDB, err := gormDB.DB()
	if err != nil {
		panic(fmt.Sprintf("postgres: failed to retrieve *sql.DB: %v", err))
	}

	return sqlDB, gormDB
}

type testingInitResource struct {
	counter atomic.Uint64
}

var (
	testInitResource *testingInitResource
	testingInitOnce  sync.Once
)

// MustInitForTest returns a function that opens a GORM *gorm.DB and the underlying standard *sql.DB using
// txdb lib (it wrap the connection inside a transaction, so it will rollback after each test)
func MustInitForTest(cfgEnv *env.Env) func(*testing.T) (*sql.DB, *gorm.DB) {
	testingInitOnce.Do(func() {
		testInitResource = &testingInitResource{}
		dsn := getDSN(cfgEnv)
		txdb.Register("txdb_pgx", "pgx", dsn)
	})

	connectFunc := func(t *testing.T) (*sql.DB, *gorm.DB) {
		db, err := sql.Open("txdb_pgx", fmt.Sprintf("%d", testInitResource.counter.Add(1)))
		if err != nil {
			panic(fmt.Sprintf("failed to open connection: %v", err))
		}

		t.Cleanup(func() {
			db.Close()
		})

		return db, initGorm(db, true)
	}

	return connectFunc
}
