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

type RegistrationHandler struct {
	svc      *service.RegistrationService
	validate *validator.Validate
}

func NewRegistrationHandler(svc *service.RegistrationService, v *validator.Validate) *RegistrationHandler {
	return &RegistrationHandler{svc: svc, validate: v}
}

func (h *RegistrationHandler) ApplyStudent(c *fiber.Ctx) error {
	var req dto.StudentRegistrationRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	reg, err := h.svc.ApplyStudent(req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(dto.ToRegistrationDTO(reg))
}

func (h *RegistrationHandler) ApplyCoach(c *fiber.Ctx) error {
	var req dto.CoachRegistrationRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	reg, err := h.svc.ApplyCoach(req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(dto.ToRegistrationDTO(reg))
}

func (h *RegistrationHandler) List(c *fiber.Ctx) error {
	p := pagination.From(c)
	kind := c.Query("kind")
	status := c.Query("status")
	search := c.Query("search")

	items, total, err := h.svc.List(kind, status, search, p.Offset(), p.Limit())
	if err != nil {
		return err
	}

	dtoItems := make([]dto.RegistrationDTO, 0, len(items))
	for i := range items {
		dtoItems = append(dtoItems, dto.ToRegistrationDTO(&items[i]))
	}
	return c.JSON(pagination.NewPage(dtoItems, total, p))
}

func (h *RegistrationHandler) Decide(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "geçersiz id")
	}
	var req dto.RegistrationDecisionRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	actorID, ok := middleware.UserID(c)
	if !ok {
		return apperrors.ErrUnauthorized
	}
	res, err := h.svc.Decide(c.Context(), id, actorID, req)
	if err != nil {
		return err
	}
	return c.JSON(res)
}
