package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/koc-luk/backend/internal/dto"
	"github.com/koc-luk/backend/internal/middleware"
	"github.com/koc-luk/backend/internal/repository"
	"github.com/koc-luk/backend/internal/service"
	apperrors "github.com/koc-luk/backend/pkg/errors"
)

type AuthHandler struct {
	auth     *service.AuthService
	users    *repository.UserRepository
	validate *validator.Validate
}

func NewAuthHandler(auth *service.AuthService, users *repository.UserRepository, v *validator.Validate) *AuthHandler {
	return &AuthHandler{auth: auth, users: users, validate: v}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	res, err := h.auth.Login(req)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID, ok := middleware.UserID(c)
	if !ok {
		return apperrors.ErrUnauthorized
	}
	u, err := h.users.FindByID(userID)
	if err != nil {
		return err
	}
	if u == nil {
		return apperrors.ErrNotFound
	}
	return c.JSON(dto.ToUserDTO(u))
}
