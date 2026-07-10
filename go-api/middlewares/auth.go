package middlewares

import (
	"strings"

	"go-api/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth valida el token JWT del header Authorization y guarda token y usuario en el contexto.
func JWTAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return c.Status(fiber.StatusUnauthorized).
				JSON(config.CustomError(config.ErrUnauthorized, "Header Authorization no encontrado"))
		}
		if !strings.HasPrefix(header, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).
				JSON(config.CustomError(config.ErrUnauthorized, "El header Authorization debe usar el esquema Bearer"))
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).
				JSON(config.CustomError(config.ErrUnauthorized, "Token inválido o expirado"))
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if sub, ok := claims["sub"].(string); ok {
				c.Locals("user", sub)
			}
		}
		c.Locals("token", tokenString)

		return c.Next()
	}
}
