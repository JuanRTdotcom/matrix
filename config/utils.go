package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// HumanReadableValidationMessage traduce un error del validador a un mensaje legible en español.
func HumanReadableValidationMessage(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())
	param := fe.Param()

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("El campo '%s' es obligatorio.", field)
	case "min":
		return fmt.Sprintf("El campo '%s' debe tener al menos %s elementos.", field, param)
	case "max":
		return fmt.Sprintf("El campo '%s' no debe exceder %s elementos.", field, param)
	case "gt":
		return fmt.Sprintf("El campo '%s' debe ser mayor que %s.", field, param)
	default:
		return fmt.Sprintf("El campo '%s' no cumple con la regla '%s'.", field, fe.Tag())
	}
}

// JsonFieldName devuelve el nombre json del campo que falló la validación.
func JsonFieldName(fe validator.FieldError, obj interface{}) string {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if field, ok := t.FieldByName(fe.StructField()); ok {
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			return strings.Split(jsonTag, ",")[0]
		}
	}
	return strings.ToLower(fe.Field())
}
