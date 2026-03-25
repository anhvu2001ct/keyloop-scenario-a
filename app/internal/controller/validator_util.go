package controller

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var _validator *validator.Validate

type errValidateRequest struct {
	FieldErrors map[string]string
}

func (e *errValidateRequest) Error() string {
	return fmt.Sprintf("validation errors: %d", len(e.FieldErrors))
}

// validateStruct validate the passing struct.
func validateStruct(value any) error {
	err := _validator.Struct(value)
	if err == nil {
		return nil
	}
	errors := err.(validator.ValidationErrors)
	fieldErrors := make(map[string]string)
	for _, err := range errors {
		key := err.Namespace()
		if key == "" {
			key = err.Field()
		} else {
			key = strings.SplitN(key, ".", 2)[1]
		}
		fieldErrors[key] = err.Tag()
	}

	return &errValidateRequest{FieldErrors: fieldErrors}
}

func init() {
	_validator = validator.New(validator.WithRequiredStructEnabled())
	_validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
