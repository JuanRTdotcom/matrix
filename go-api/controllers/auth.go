package controllers

import (
	"errors"

	"go-api/config"
	"go-api/models"
	"go-api/usecases"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// AuthHandler expone los endpoints de autenticación.
type AuthHandler struct {
	usecase  Authenticator
	validate *validator.Validate
}

// NewAuthHandler construye el handler de autenticación.
func NewAuthHandler(u Authenticator) *AuthHandler {
	return &AuthHandler{usecase: u, validate: validator.New()}
}

// Login valida credenciales y devuelve un token JWT.
// @Summary Autenticación y emisión de token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.LoginRequest true "Credenciales"
// @Success 200 {object} models.LoginResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(config.CustomError(config.ErrBadRequest, "JSON inválido"))
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(buildValidationError(err, req))
	}

	res, err := h.usecase.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, usecases.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).
				JSON(config.CustomError(config.ErrUnauthorized, "Usuario o contraseña incorrectos"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(config.ErrInternal)
	}

	return c.Status(fiber.StatusOK).JSON(config.SuccessOK("Autenticación exitosa", res))
}

// buildValidationError transforma los errores del validador en la respuesta estándar.
func buildValidationError(err error, obj interface{}) config.APIError {
	var fieldErrors []config.FieldError
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		for _, fe := range verrs {
			fieldErrors = append(fieldErrors, config.FieldError{
				Field:   config.JsonFieldName(fe, obj),
				Message: config.HumanReadableValidationMessage(fe),
			})
		}
	}
	return config.CustomErrors(fieldErrors)
}
