package config

// APISuccess es la estructura estándar de respuesta exitosa de la API.
type APISuccess struct {
	Code    int    `json:"code" example:"200"`
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Operación exitosa"`
	Data    any    `json:"data,omitempty"`
}

// SuccessOK construye una respuesta 200 con mensaje y data.
func SuccessOK(message string, data any) APISuccess {
	return APISuccess{
		Code:    200,
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

// SuccessCreated construye una respuesta 201 con mensaje y data.
func SuccessCreated(message string, data any) APISuccess {
	return APISuccess{
		Code:    201,
		Status:  "success",
		Message: message,
		Data:    data,
	}
}
