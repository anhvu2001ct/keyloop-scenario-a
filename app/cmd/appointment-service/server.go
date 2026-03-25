package main

import (
	"context"
	"fmt"
	"scenario-a/internal/dep"
	"scenario-a/internal/route"
	"scenario-a/pkg/logger"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func initRouter(deps *dep.Dependencies) *echo.Echo {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestID())
	e.Use(echo.WrapMiddleware(otelhttp.NewMiddleware("appointment-service")))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		BeforeNextFunc: func(c *echo.Context) {
			fields := []zap.Field{
				zap.String("method", c.Request().Method),
				zap.String("uri", c.Request().RequestURI),
				zap.String("requestID", c.Response().Header().Get(echo.HeaderXRequestID)),
			}

			spanCtx := trace.SpanContextFromContext(c.Request().Context())
			if spanCtx.IsValid() {
				fields = append(fields,
					zap.String("trace_id", spanCtx.TraceID().String()),
					zap.String("span_id", spanCtx.SpanID().String()),
				)
			}

			logger.Info("Incoming Request", fields...)
		},
		LogMethod:    true,
		LogURI:       true,
		LogStatus:    true,
		LogRequestID: true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			fields := []zap.Field{
				zap.String("method", v.Method),
				zap.String("uri", v.URI),
				zap.String("requestID", v.RequestID),
				zap.Int("responseCode", v.Status),
			}

			spanCtx := trace.SpanContextFromContext(c.Request().Context())
			if spanCtx.IsValid() {
				fields = append(fields,
					zap.String("trace_id", spanCtx.TraceID().String()),
				)
			}

			logger.Info("Response Request", fields...)
			return nil
		},
	}))
	e.Use(middleware.Recover())

	route.Load(e, deps)
	return e
}

func startServer(ctx context.Context, port int, deps *dep.Dependencies) {
	router := initRouter(deps)

	sc := echo.StartConfig{
		Address:         fmt.Sprintf(":%d", port),
		GracefulTimeout: 5 * time.Second,
	}
	if err := sc.Start(ctx, router); err != nil {
		router.Logger.Error("failed to start server", "error", err)
	}
}
