package middlewares

import (
	"go-api/config"

	"github.com/gofiber/fiber/v2"
)

// NotFoundMiddleware maneja las rutas no registradas devolviendo un 404 estándar.
func NotFoundMiddleware(c *fiber.Ctx) error {
	return c.Status(config.ErrNotFound.Code).JSON(config.ErrNotFound)
}
