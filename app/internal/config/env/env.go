package env

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"

	"github.com/joho/godotenv"
)

var (
	ErrGetCallerInfo = errors.New("cannot retrieve caller info")
)

type EnvMode string

func (mode EnvMode) Valid() bool {
	return slices.Contains([]EnvMode{Development, Test, Production}, mode)
}

var (
	Development EnvMode = "development"
	Test        EnvMode = "test"
	Production  EnvMode = "production"
)

type Env struct {
	AppEnv    EnvMode
	AppPort   int
	LogFormat string

	// Database
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func MustLoad(mode EnvMode) *Env {
	if !mode.Valid() {
		panic("invalid env value: " + mode)
	}

	if err := loadDotEnv(".env", false, true); err != nil {
		panic(fmt.Sprintf("cannot load .env: %v", err))
	}

	secondEnvFile := ".env." + string(mode)

	if err := loadDotEnv(secondEnvFile, true, true); err != nil {
		panic(fmt.Sprintf("cannot load %s: %v", secondEnvFile, err))
	}

	return &Env{
		AppEnv:    EnvMode(mode),
		AppPort:   Int("APP_PORT"),
		LogFormat: String("LOG_FORMAT", true),

		DBHost:     String("DB_HOST", false),
		DBPort:     Int("DB_PORT"),
		DBUser:     String("DB_USER", false),
		DBPassword: String("DB_PASSWORD", false),
		DBName:     String("DB_NAME", false),
		DBSSLMode:  String("DB_SSLMODE", false),
	}
}

func loadDotEnv(envFileName string, overwrite, ignoreIfNotFound bool) error {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ErrGetCallerInfo
	}

	file = filepath.Join(file, "../../../..", envFileName)

	var err error
	if overwrite {
		err = godotenv.Overload(file)
	} else {
		err = godotenv.Load(file)
	}

	if ignoreIfNotFound && errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}
