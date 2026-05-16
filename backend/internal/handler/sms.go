package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/koc-luk/backend/internal/dto"
	"github.com/koc-luk/backend/internal/middleware"
	"github.com/koc-luk/backend/internal/service"
	apperrors "github.com/koc-luk/backend/pkg/errors"
	"github.com/koc-luk/backend/pkg/pagination"
)

type SmsHandler struct {
	svc      *service.SmsService
	validate *validator.Validate
}

func NewSmsHandler(svc *service.SmsService, v *validator.Validate) *SmsHandler {
	return &SmsHandler{svc: svc, validate: v}
}

func (h *SmsHandler) Templates(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"items": h.svc.Templates()})
}

func (h *SmsHandler) Send(c *fiber.Ctx) error {
	uid, _ := middleware.UserID(c)
	var req dto.SendSmsRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	res, err := h.svc.Send(c.Context(), uid, req)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (h *SmsHandler) Logs(c *fiber.Ctx) error {
	p := pagination.From(c)
	items, total, err := h.svc.Logs(c.Query("search"), p.Offset(), p.Limit())
	if err != nil {
		return err
	}
	return c.JSON(pagination.NewPage(items, total, p))
}

func (h *SmsHandler) MissingCoaches(c *fiber.Ctx) error {
	weekID := c.Query("weekId")
	if weekID == "" {
		return apperrors.WithDetails(apperrors.ErrValidation, "weekId gerekli")
	}
	items, err := h.svc.MissingCoachesForWeek(weekID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}
