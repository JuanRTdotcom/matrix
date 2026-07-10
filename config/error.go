package config

// APIError es la estructura estándar de error de la API.
type APIError struct {
	Code    int          `json:"code" example:"404"`
	Status  bool         `json:"status" example:"false"`
	Error   string       `json:"error" example:"not_found"`
	Message string       `json:"message" example:"El recurso solicitado no existe"`
	Errors  []FieldError `json:"errors,omitempty"`
}

// FieldError describe el error de un campo específico en la validación.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Errores del cliente (4xx)
var (
	ErrBadRequest = APIError{
		Code:    400,
		Status:  false,
		Error:   "bad_request",
		Message: "Los campos enviados no son válidos",
	}

	ErrUnauthorized = APIError{
		Code:    401,
		Status:  false,
		Error:   "unauthorized",
		Message: "No autorizado para acceder a este recurso",
	}

	ErrNotFound = APIError{
		Code:    404,
		Status:  false,
		Error:   "not_found",
		Message: "El recurso solicitado no fue encontrado",
	}

	ErrUnprocessableEntity = APIError{
		Code:    422,
		Status:  false,
		Error:   "unprocessable_entity",
		Message: "No se pudo procesar la solicitud por errores de validación",
	}
)

// Errores del servidor (5xx)
var (
	ErrInternal = APIError{
		Code:    500,
		Status:  false,
		Error:   "internal_error",
		Message: "Error interno del servidor",
	}

	ErrServiceUnavailable = APIError{
		Code:    503,
		Status:  false,
		Error:   "service_unavailable",
		Message: "Servicio no disponible, intente más tarde",
	}
)

// CustomError reutiliza un error base con un mensaje personalizado.
func CustomError(base APIError, message string) APIError {
	return APIError{
		Code:    base.Code,
		Status:  base.Status,
		Error:   base.Error,
		Message: message,
	}
}

// CustomErrors adjunta errores de campo al error de bad request.
func CustomErrors(fieldErrors []FieldError) APIError {
	return APIError{
		Code:    ErrBadRequest.Code,
		Status:  ErrBadRequest.Status,
		Error:   ErrBadRequest.Error,
		Message: ErrBadRequest.Message,
		Errors:  fieldErrors,
	}
}
