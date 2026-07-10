package controllers

import "github.com/gofiber/fiber/v2"

// Bienvenido responde el estado y metadatos de la API.
// @Summary Endpoint de bienvenida / healthcheck
// @Tags General
// @Produce json
// @Router / [get]
func Bienvenido(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  true,
		"code":    200,
		"message": "API en Go (Fiber) — Factorización QR + estadísticas (vía Node.js)",
		"version": "1.0.0",
	})
}
