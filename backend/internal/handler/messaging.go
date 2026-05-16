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

type MessagingHandler struct {
	svc      *service.MessagingService
	validate *validator.Validate
}

func NewMessagingHandler(svc *service.MessagingService, v *validator.Validate) *MessagingHandler {
	return &MessagingHandler{svc: svc, validate: v}
}

func (h *MessagingHandler) ListThreads(c *fiber.Ctx) error {
	uid, _ := middleware.UserID(c)
	role, _ := middleware.Role(c)
	p := pagination.From(c)
	items, total, err := h.svc.ListThreads(uid, role, c.Query("kind"), c.Query("status"), c.Query("search"), p.Offset(), p.Limit())
	if err != nil {
		return err
	}
	return c.JSON(pagination.NewPage(items, total, p))
}

func (h *MessagingHandler) Create(c *fiber.Ctx) error {
	uid, _ := middleware.UserID(c)
	role, _ := middleware.Role(c)
	var req dto.CreateThreadRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	d, err := h.svc.Create(uid, role, req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(d)
}

func (h *MessagingHandler) Messages(c *fiber.Ctx) error {
	uid, _ := middleware.UserID(c)
	role, _ := middleware.Role(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "geçersiz id")
	}
	items, err := h.svc.Messages(id, uid, role)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *MessagingHandler) Send(c *fiber.Ctx) error {
	uid, _ := middleware.UserID(c)
	role, _ := middleware.Role(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "geçersiz id")
	}
	var req dto.SendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	d, err := h.svc.Send(id, uid, role, req.Body)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(d)
}

func (h *MessagingHandler) Close(c *fiber.Ctx) error {
	uid, _ := middleware.UserID(c)
	role, _ := middleware.Role(c)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "geçersiz id")
	}
	if err := h.svc.Close(id, uid, role); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"ok": true})
}
