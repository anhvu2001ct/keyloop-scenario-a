package config

import (
	"scenario-a/internal/config/env"
	"scenario-a/pkg/logger"
	"sync"
)

type Config struct {
	Env *env.Env
}

var initOnce sync.Once
var config *Config

func Get() *Config {
	return config
}

func MustInit() {
	mustInit(false)
}

func MustInitForTest() {
	mustInit(true)
}

func mustInit(forTest bool) {
	initOnce.Do(func() {
		envMode := env.Test
		if !forTest {
			envMode = env.EnvMode(env.String("APP_ENV", true))
		}

		cfgEnv := env.MustLoad(envMode)

		logger.MustInit(string(envMode), cfgEnv.LogFormat)

		config = &Config{
			Env: cfgEnv,
		}
	})
}
