package usecases

import (
	"errors"
	"time"

	"go-api/models"

	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalidCredentials se devuelve cuando el usuario o contraseña no son válidos.
var ErrInvalidCredentials = errors.New("credenciales inválidas")

// AuthUseCase encapsula la emisión de tokens JWT. Para el reto las credenciales
// válidas se configuran por variables de entorno (usuario/contraseña de demo).
type AuthUseCase struct {
	jwtSecret   []byte
	validUser   string
	validPass   string
	tokenExpiry time.Duration
}

// NewAuthUseCase construye el caso de uso de autenticación.
func NewAuthUseCase(jwtSecret, validUser, validPass string, expiry time.Duration) *AuthUseCase {
	return &AuthUseCase{
		jwtSecret:   []byte(jwtSecret),
		validUser:   validUser,
		validPass:   validPass,
		tokenExpiry: expiry,
	}
}

// Login valida las credenciales y, si son correctas, emite un token JWT firmado.
func (uc *AuthUseCase) Login(username, password string) (models.LoginResponse, error) {
	if username != uc.validUser || password != uc.validPass {
		return models.LoginResponse{}, ErrInvalidCredentials
	}

	now := time.Now()
	expiresAt := now.Add(uc.tokenExpiry)

	claims := jwt.MapClaims{
		"sub": username,
		"iat": now.Unix(),
		"exp": expiresAt.Unix(),
		"iss": "go-api",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(uc.jwtSecret)
	if err != nil {
		return models.LoginResponse{}, err
	}

	return models.LoginResponse{
		Token:     signed,
		ExpiresIn: int64(uc.tokenExpiry.Seconds()),
		Type:      "Bearer",
	}, nil
}
