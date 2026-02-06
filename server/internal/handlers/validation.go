package handlers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Глобальный валидатор — создаётся один раз
var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type ValidationError struct {
	Error     string `json:error`
	Field     string `json: field`
	Condition string `json: condition`
}

func (h *Handler) ValidateBodyRequest(payload interface{}) []*ValidationError {
	var errors []*ValidationError

	err := validate.Struct(payload) // ← используем глобальный validate
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			reflected := reflect.ValueOf(payload)

			for _, validationError := range validationErrors {
				field, _ := reflected.Type().FieldByName(validationError.StructField())
				jsonField := field.Tag.Get("json")
				if jsonField == "" || jsonField == "-" {
					jsonField = strings.ToLower(validationError.StructField())
				}

				condition := validationError.Tag()
				param := validationError.Param() // параметр проверки, если есть

				errorMessage := fmt.Sprintf("'%s' haven't pass validation '%s'", jsonField, condition)
				if param != "" {
					errorMessage += fmt.Sprintf(" (%s)", param)
				}

				switch condition {
				case "eqfield":
					paramField, _ := reflected.Type().FieldByName(param)
					jsonParamField := paramField.Tag.Get("json")
					if jsonParamField == "" || jsonParamField == "-" {
						jsonParamField = strings.ToLower(validationError.StructField())
					}
					errorMessage = fmt.Sprintf("Field '%s' must be equal to field '%s'", jsonField, jsonParamField)
				}

				currentValidationError := ValidationError{
					Error:     errorMessage,
					Field:     jsonField,
					Condition: condition,
				}
				errors = append(errors, &currentValidationError)
			}
		}
	}
	return errors
}
