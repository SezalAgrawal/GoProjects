package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

var (
	UnexpectedErr  = "Unexpected"
	ErrInvalidType = func(field string, expectedType interface{}, err error) ValidationErrorInterface {
		return &ValidationError{
			errorType: "InvalidType",
			message:   fmt.Sprintf("InvalidType for field: %s. Expected: %s", field, expectedType),
			Err:       err,
		}
	}
	ErrInvalidJson = func(err error) ValidationErrorInterface {
		return &ValidationError{
			errorType: "InvalidJson",
			message:   fmt.Sprintf("InvalidJson: %s", err.Error()),
			Err:       err,
		}
	}
	ErrInvalidValue = func(message string, err error) ValidationErrorInterface {
		if message == "" {
			message = err.Error()
		}
		return &ValidationError{
			errorType: "InvalidValue",
			message:   fmt.Sprintf("InvalidValue: %s", message),
			Err:       err,
		}
	}
)

type ValidationErrorInterface interface {
	Type() string
	Error() string
	Unwrap() error
	IsUnexpectedErr() bool
}

type ValidationError struct {
	errorType string
	message   string
	Err       error
}

func (e *ValidationError) Error() string { return e.message }

func (e *ValidationError) Type() string { return e.errorType }

func (e *ValidationError) Unwrap() error { return e.Err }

func (e *ValidationError) IsUnexpectedErr() bool { return e.errorType == UnexpectedErr }

func ValidateMapToStruct(m map[string]interface{}, s interface{}, structValidations ...validator.StructLevelFunc) ValidationErrorInterface {
	jsonString, _ := json.Marshal(m)
	err := json.Unmarshal(jsonString, &s)
	if err != nil {
		return HandleValidationErrors(err)
	}
	return ValidateStruct(s, structValidations...)
}

func ValidateStruct(s interface{}, structValidations ...validator.StructLevelFunc) ValidationErrorInterface {
	var validate = validator.New()
	_ = validate.RegisterValidation("notblank", validators.NotBlank)
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	for _, sValidation := range structValidations {
		validate.RegisterStructValidation(sValidation, s)
	}

	err := validate.Struct(s)
	if err != nil {
		return HandleValidationErrors(err)
	}
	return nil
}

func HandleValidationErrors(err error) ValidationErrorInterface {
	switch e := err.(type) {
	case *json.UnmarshalTypeError:
		return ErrInvalidType(e.Field, e.Type, e)
	case validator.ValidationErrors:
		var msgs []string
		for _, fe := range e {
			namespace := fe.Namespace()
			namespaces := strings.Split(namespace, ".")

			filteredNamespaces := make([]string, 0)
			for _, ns := range namespaces {
				// remove the struct name prefix, the structs are mixedCase the tags are lower case always
				//	eg1: orderRequest.customer.id =>  ["customer", "id"]
				//  eg2: orderCreateRequest.baseOrder.customer.id => ["customer", "id"]
				if ns != strings.ToLower(ns) {
					continue
				}
				filteredNamespaces = append(filteredNamespaces, ns)
			}
			fieldName := strings.Join(filteredNamespaces, ".")
			if fieldName == "" {
				fieldName = fe.Field()
			}
			switch fe.Tag() {
			case "required":
				msgs = append(msgs, fmt.Sprintf("%s is a required field", fieldName))
			case "notblank":
				msgs = append(msgs, fmt.Sprintf("%s should not be empty", fieldName))
			case "max":
				msgs = append(msgs, fmt.Sprintf("%s must be a maximum of %s in length", fieldName, fe.Param()))
			case "url":
				msgs = append(msgs, fmt.Sprintf("%s must be a valid URL", fieldName))
			case "uuid":
				msgs = append(msgs, fmt.Sprintf("%s must be a valid uuid", fieldName))
			case "oneof":
				msgs = append(msgs, fmt.Sprintf("%s value must be one of the pre-decided values", fieldName))
			case "gt":
				msgs = append(msgs, fmt.Sprintf("%s must be a greater than %s", fieldName, fe.Param()))
			case "gte":
				msgs = append(msgs, fmt.Sprintf("%s must be a greater than or equal to %s", fieldName, fe.Param()))
			case "lt":
				msgs = append(msgs, fmt.Sprintf("%s must be a less than %s", fieldName, fe.Param()))
			case "lte":
				msgs = append(msgs, fmt.Sprintf("%s must be a less than or equal to %s", fieldName, fe.Param()))
			default:
				msgs = append(msgs, fmt.Sprintf("%s:%s", fieldName, fe.Tag()))
			}
		}
		return ErrInvalidValue(strings.Join(msgs, ", "), e)
	case *json.SyntaxError:
		return ErrInvalidJson(e)
	default:
		return &ValidationError{
			errorType: UnexpectedErr,
			message:   e.Error(),
			Err:       e,
		}
	}
}
