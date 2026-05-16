package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/koc-luk/backend/internal/dto"
	"github.com/koc-luk/backend/internal/middleware"
	"github.com/koc-luk/backend/internal/service"
	apperrors "github.com/koc-luk/backend/pkg/errors"
)

type WeekHandler struct {
	svc      *service.WeekService
	validate *validator.Validate
}

func NewWeekHandler(svc *service.WeekService, v *validator.Validate) *WeekHandler {
	return &WeekHandler{svc: svc, validate: v}
}

func (h *WeekHandler) ListAll(c *fiber.Ctx) error {
	ws, err := h.svc.ListAll()
	if err != nil {
		return err
	}
	items := make([]dto.WeekDTO, 0, len(ws))
	for i := range ws {
		items = append(items, dto.ToWeekDTO(&ws[i]))
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *WeekHandler) ListOpen(c *fiber.Ctx) error {
	ws, err := h.svc.ListOpen()
	if err != nil {
		return err
	}
	items := make([]dto.WeekDTO, 0, len(ws))
	for i := range ws {
		items = append(items, dto.ToWeekDTO(&ws[i]))
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *WeekHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateWeekRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	w, err := h.svc.Create(req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(dto.ToWeekDTO(w))
}

func (h *WeekHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "geçersiz id")
	}
	var req dto.UpdateWeekRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	w, err := h.svc.Update(id, req)
	if err != nil {
		return err
	}
	return c.JSON(dto.ToWeekDTO(w))
}

func (h *WeekHandler) Open(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "geçersiz id")
	}
	by, ok := middleware.UserID(c)
	if !ok {
		return apperrors.ErrUnauthorized
	}
	w, err := h.svc.Open(id, by)
	if err != nil {
		return err
	}
	return c.JSON(dto.ToWeekDTO(w))
}

func (h *WeekHandler) Close(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "geçersiz id")
	}
	w, err := h.svc.Close(id)
	if err != nil {
		return err
	}
	return c.JSON(dto.ToWeekDTO(w))
}

func (h *WeekHandler) Generate(c *fiber.Ctx) error {
	var req dto.GenerateWeeksRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	created, err := h.svc.Generate(req)
	if err != nil {
		return err
	}
	items := make([]dto.WeekDTO, 0, len(created))
	for i := range created {
		items = append(items, dto.ToWeekDTO(&created[i]))
	}
	return c.JSON(fiber.Map{"items": items, "created": len(items)})
}
