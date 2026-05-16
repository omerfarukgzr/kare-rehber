package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/koc-luk/backend/internal/dto"
	"github.com/koc-luk/backend/internal/middleware"
	"github.com/koc-luk/backend/internal/service"
	apperrors "github.com/koc-luk/backend/pkg/errors"
	"github.com/koc-luk/backend/pkg/pagination"
)

type UserHandler struct {
	svc      *service.UserService
	validate *validator.Validate
}

func NewUserHandler(svc *service.UserService, v *validator.Validate) *UserHandler {
	return &UserHandler{svc: svc, validate: v}
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	p := pagination.From(c)
	role := c.Query("role")
	isActive := c.Query("isActive")
	city := c.Query("city")
	search := c.Query("search")

	users, total, err := h.svc.List(role, isActive, city, search, p.Offset(), p.Limit())
	if err != nil {
		return err
	}
	items := make([]dto.UserDTO, 0, len(users))
	for i := range users {
		items = append(items, dto.ToUserDTO(&users[i]))
	}
	return c.JSON(pagination.NewPage(items, total, p))
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "geçersiz id")
	}
	u, err := h.svc.Get(id)
	if err != nil {
		return err
	}
	return c.JSON(dto.ToUserDTO(u))
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	actorID, ok := middleware.UserID(c)
	if !ok {
		return apperrors.ErrUnauthorized
	}
	u, plain, err := h.svc.Create(c.Context(), actorID, req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(dto.CreateUserResponse{
		User:              dto.ToUserDTO(u),
		GeneratedPassword: plain,
	})
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "geçersiz id")
	}
	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	u, err := h.svc.Update(id, req)
	if err != nil {
		return err
	}
	return c.JSON(dto.ToUserDTO(u))
}

func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "geçersiz id")
	}
	actorID, ok := middleware.UserID(c)
	if !ok {
		return apperrors.ErrUnauthorized
	}
	plain, err := h.svc.ResetPassword(c.Context(), actorID, id)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"generatedPassword": plain})
}
