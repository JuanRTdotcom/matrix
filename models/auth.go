package models

// LoginRequest son las credenciales de acceso.
// @Description Credenciales para obtener un token JWT.
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse contiene el token JWT emitido.
// @Description Token de acceso.
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expiresIn"`
	Type      string `json:"type"`
}
