package controller

import (
	"fmt"
	"scenario-a/pkg/common"

	"github.com/labstack/echo/v5"
)

func bindBody(c *echo.Context, req any) error {
	if err := c.Bind(req); err != nil {
		return &common.AppError{
			Type:    common.AppErrorValidationFailed,
			Message: "invalid field(s) format or data type",
		}
	}
	return nil
}

func getPathParam[T any](c *echo.Context, key string) (T, error) {
	val, err := echo.PathParam[T](c, key)
	if err != nil {
		return val, &common.AppError{
			Type:    common.AppErrorValidationFailed,
			Message: fmt.Sprintf("invalid path param '%s', cannot convert to '%T'", key, val),
		}
	}
	return val, nil
}
