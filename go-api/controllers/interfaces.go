package controllers

import "go-api/models"

// Interfaces que los controladores usan, para no depender de implementaciones concretas.

// Authenticator abstrae el caso de uso de autenticación.
type Authenticator interface {
	Login(username, password string) (models.LoginResponse, error)
}

// Processor abstrae el caso de uso de procesamiento de matrices.
type Processor interface {
	Execute(token string, matrix [][]float64) (models.ProcessResponse, error)
}

// ClientError marca los errores de entrada inválida, que el controlador mapea a 422.
type ClientError interface {
	error
	IsClientError() bool
}
