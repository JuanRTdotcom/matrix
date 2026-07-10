package usecases

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuth_LoginSuccessIssuesToken(t *testing.T) {
	secret := "test-secret"
	uc := NewAuthUseCase(secret, "admin", "admin123", time.Hour)

	res, err := uc.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if res.Token == "" {
		t.Fatal("se esperaba un token no vacío")
	}
	if res.Type != "Bearer" {
		t.Errorf("tipo = %q, se esperaba Bearer", res.Type)
	}

	// El token debe ser válido y firmado con el mismo secreto, con el sub correcto.
	token, err := jwt.Parse(res.Token, func(_ *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		t.Fatalf("token inválido: %v", err)
	}
	claims := token.Claims.(jwt.MapClaims)
	if claims["sub"] != "admin" {
		t.Errorf("sub = %v, se esperaba admin", claims["sub"])
	}
}

func TestAuth_LoginRejectsBadCredentials(t *testing.T) {
	uc := NewAuthUseCase("secret", "admin", "admin123", time.Hour)

	if _, err := uc.Login("admin", "wrong"); !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("se esperaba ErrInvalidCredentials, se obtuvo %v", err)
	}
	if _, err := uc.Login("nobody", "admin123"); !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("se esperaba ErrInvalidCredentials, se obtuvo %v", err)
	}
}
