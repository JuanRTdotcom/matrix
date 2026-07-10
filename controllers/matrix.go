package controllers

import (
	"errors"

	"go-api/config"
	"go-api/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// MatrixHandler expone el endpoint de procesamiento de matrices.
type MatrixHandler struct {
	usecase  Processor
	validate *validator.Validate
}

// NewMatrixHandler construye el handler de matrices.
func NewMatrixHandler(u Processor) *MatrixHandler {
	return &MatrixHandler{usecase: u, validate: validator.New()}
}

// Process factoriza la matriz (QR) y devuelve el QR junto a sus estadísticas.
// @Summary Procesa una matriz (QR + estadísticas)
// @Tags Matriz
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.MatrixRequest true "Matriz de entrada"
// @Success 200 {object} models.ProcessResponse
// @Router /api/v1/matrix/process [post]
func (h *MatrixHandler) Process(c *fiber.Ctx) error {
	var req models.MatrixRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(config.CustomError(config.ErrBadRequest, "JSON inválido"))
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(buildValidationError(err, req))
	}

	// El token ya lo validó el middleware; se reenvía a Node.
	token, _ := c.Locals("token").(string)

	res, err := h.usecase.Execute(token, req.Matrix)
	if err != nil {
		// Entrada inválida del cliente -> 422.
		var clientErr ClientError
		if errors.As(err, &clientErr) && clientErr.IsClientError() {
			return c.Status(fiber.StatusUnprocessableEntity).
				JSON(config.CustomError(config.ErrUnprocessableEntity, clientErr.Error()))
		}
		// Fallo de la API de Node -> 503.
		return c.Status(fiber.StatusServiceUnavailable).
			JSON(config.CustomError(config.ErrServiceUnavailable, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(config.SuccessOK("Matriz procesada correctamente", res))
}
