package server

import (
	"go-api/config"
	"go-api/controllers"
	"go-api/middlewares"

	"github.com/gofiber/fiber/v2"
)

// registerRoutes define todas las rutas de la API en un solo lugar.
func registerRoutes(app *fiber.App, cfg *config.Config, auth *controllers.AuthHandler, matrix *controllers.MatrixHandler) {
	// Público (sin JWT).
	app.Get("/", controllers.Bienvenido)       // Healthcheck
	app.Post("/api/v1/auth/login", auth.Login) // Login: devuelve un JWT

	// Protegido (requiere JWT).
	api := app.Group("/api/v1", middlewares.JWTAuth(cfg.JWTSecret))
	api.Post("/matrix/process", matrix.Process) // Factoriza (QR) y devuelve QR + estadísticas

	app.Use(middlewares.NotFoundMiddleware) // 404 para rutas no registradas
}
