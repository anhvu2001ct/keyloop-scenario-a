package sqlmigrate

import (
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"scenario-a/internal/config/env"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
)

type SqlMigration struct {
	db *dbmate.DB
}

func New(envMode env.EnvMode, engine string) *SqlMigration {
	dsn := buildDSN(envMode, engine)
	url, err := url.Parse(dsn)
	if err != nil {
		panic("error while parsing dsn: " + err.Error())
	}

	db := dbmate.New(url)
	db.AutoDumpSchema = false
	db.MigrationsDir = []string{getMigrationFolderPath()}
	db.SchemaFile = filepath.Join(db.MigrationsDir[0], "..", "schema.sql")

	return &SqlMigration{db}
}

func buildDSN(envMode env.EnvMode, engine string) string {
	protocol := "postgres"
	if engine == "mysql" {
		protocol = "mysql"
	}

	envCfg := env.MustLoad(envMode)
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s", protocol, envCfg.DBUser, envCfg.DBPassword, envCfg.DBHost, envCfg.DBPort, envCfg.DBName, envCfg.DBSSLMode)
}

func getMigrationFolderPath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot retrieve file's caller info")
	}

	rootPath := filepath.Join(filepath.Dir(file), "../..")
	return filepath.Join(rootPath, "db", "migrations")
}
