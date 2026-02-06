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
				param := validationError.Param() // для oneof это будет "DEPOSIT WITHDRAW"

				errorMessage := fmt.Sprintf("'%s' haven't pass validation '%s'", jsonField, condition)
				if param != "" {
					errorMessage += fmt.Sprintf(" (%s)", param)
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
