package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type ValidationError struct {
	Field   string
	Message string
	Tag     string
}

func init() {
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func MapValidationError(err error) []ValidationError {
	var errs []ValidationError
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		for _, fe := range ve {
			errs = append(errs, ValidationError{
				Field:   fe.Field(),
				Message: getErrorMsg(fe),
				Tag:     fe.ActualTag(),
			})
		}
	}
	return errs
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Email is invalid"
	case "min":
		return fmt.Sprintf("Minimum length is %s characters", fe.Param())
	}
	return "Value is invalid"
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
