package util

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   any    `json:"value"`
}

var (
	Validator *validator.Validate
)

func Validate[T any](dto T) []ValidationError {
	Validator = validator.New(validator.WithRequiredStructEnabled())
	validationErrors := []ValidationError{}
	err := Validator.Struct(dto)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.StructField()
			field, _ := reflect.TypeOf(dto).FieldByName(fieldName)
			tagValue := field.Tag.Get("json")
			if tagValue == "" {
				tagValue = err.Field()
			}
			validationError := ValidationError{
				Field:   tagValue,
				Tag:     err.Tag(),
				Value:   err.Value(),
				Message: err.Error(),
			}
			validationErrors = append(validationErrors, validationError)
		}
	}
	return validationErrors
}
