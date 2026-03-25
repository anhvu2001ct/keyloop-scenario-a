package controller

import (
	"errors"
	"fmt"
	"net/http"
	"scenario-a/internal/errors/dberr"
	"scenario-a/pkg/common"
	"scenario-a/pkg/logger"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

func sendError(ctx *echo.Context, err error) error {
	send := func(code int, err error) error {
		// only log error if status code >= 500
		if code >= http.StatusInternalServerError {
			if safeErr, ok := errors.AsType[*common.SafeError](err); ok {
				logger.Error(
					fmt.Sprintf("error: %s, short=%s", safeErr.UserMsg, safeErr.Short),
					zap.Int("responseCode", code),
					zap.NamedError("cause", safeErr.Cause),
					zap.Any("metadata", safeErr.Metadata),
				)
			} else {
				logger.Error(fmt.Sprintf("error: %s", err), zap.Int("responseCode", code))
			}
		}

		if validateErr, ok := errors.AsType[*errValidateRequest](err); ok {
			return ctx.JSON(code, map[string]any{
				"message": validateErr.Error(),
				"details": validateErr.FieldErrors,
			})
		}

		isSafeMessage := false
		if _, ok := errors.AsType[*common.AppError](err); ok {
			isSafeMessage = true
		} else if _, ok := errors.AsType[*common.SafeError](err); ok {
			isSafeMessage = true
		}

		errMessage := "Internal Server Error"
		// return error's message if code < 500 or this error is from a known safe error
		if isSafeMessage || code < http.StatusInternalServerError {
			errMessage = err.Error()
		}

		return ctx.JSON(code, map[string]string{
			"message": errMessage,
		})
	}

	code := http.StatusInternalServerError

	if appErr, ok := errors.AsType[*common.AppError](err); ok {
		switch appErr.Type {
		case common.AppErrorRecordNotFound,
			common.AppErrorValidationFailed,
			common.AppErrorRecordExisted:

			code = http.StatusBadRequest
		default:
			code = http.StatusInternalServerError
		}
	} else if _, ok := errors.AsType[*errValidateRequest](err); ok {
		code = http.StatusBadRequest
	} else if _, ok := errors.AsType[*dberr.RecordNotFound](err); ok {
		code = http.StatusBadRequest
	}

	return send(code, err)
}
