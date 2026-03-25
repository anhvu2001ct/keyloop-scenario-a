package main

import (
	"context"
	"os/signal"
	"scenario-a/internal/config"
	"scenario-a/internal/config/sqldb"
	"scenario-a/internal/dep"
	"scenario-a/internal/telemetry"
	"syscall"
)

func main() {
	config.MustInit()
	cfg := config.Get()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	shutdown, err := telemetry.InitTracerProvider(ctx, "appointment-service")
	if err != nil {
		panic(err)
	}
	defer shutdown(context.Background())

	_, gormDB := sqldb.MustInit(cfg.Env)
	deps := dep.Init(gormDB)

	startServer(ctx, cfg.Env.AppPort, deps)
}
