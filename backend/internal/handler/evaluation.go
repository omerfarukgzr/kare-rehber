package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/koc-luk/backend/internal/domain"
	"github.com/koc-luk/backend/internal/dto"
	"github.com/koc-luk/backend/internal/middleware"
	"github.com/koc-luk/backend/internal/service"
	apperrors "github.com/koc-luk/backend/pkg/errors"
	"github.com/koc-luk/backend/pkg/pagination"
)

type EvaluationHandler struct {
	svc      *service.EvaluationService
	coach    *service.CoachService
	validate *validator.Validate
}

func NewEvaluationHandler(svc *service.EvaluationService, coach *service.CoachService, v *validator.Validate) *EvaluationHandler {
	return &EvaluationHandler{svc: svc, coach: coach, validate: v}
}

func (h *EvaluationHandler) List(c *fiber.Ctx) error {
	uid, _ := middleware.UserID(c)
	role, _ := middleware.Role(c)
	p := pagination.From(c)

	q := service.ListQuery{
		WeekID:    c.Query("weekId"),
		CoachID:   c.Query("coachId"),
		StudentID: c.Query("studentId"),
		Status:    c.Query("status"),
		Search:    c.Query("search"),
		Offset:    p.Offset(),
		Limit:     p.Limit(),
	}
	items, total, err := h.svc.List(uid, role, q)
	if err != nil {
		return err
	}
	return c.JSON(pagination.NewPage(items, total, p))
}

func (h *EvaluationHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "geçersiz id")
	}
	role, _ := middleware.Role(c)
	d, err := h.svc.GetByID(id, role)
	if err != nil {
		return err
	}
	return c.JSON(d)
}

func (h *EvaluationHandler) CoachCreate(c *fiber.Ctx) error {
	coachID, ok := middleware.UserID(c)
	if !ok {
		return apperrors.ErrUnauthorized
	}
	var req dto.CreateEvaluationRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	d, err := h.svc.CreateForCoach(coachID, req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(d)
}

func (h *EvaluationHandler) CoachUpdate(c *fiber.Ctx) error {
	coachID, ok := middleware.UserID(c)
	if !ok {
		return apperrors.ErrUnauthorized
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "geçersiz id")
	}
	var req dto.UpdateEvaluationRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	d, err := h.svc.UpdateForCoach(coachID, id, req)
	if err != nil {
		return err
	}
	return c.JSON(d)
}

func (h *EvaluationHandler) AdminUpdate(c *fiber.Ctx) error {
	adminID, ok := middleware.UserID(c)
	if !ok {
		return apperrors.ErrUnauthorized
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "geçersiz id")
	}
	var req dto.UpdateEvaluationRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	d, err := h.svc.UpdateByAdmin(adminID, id, req)
	if err != nil {
		return err
	}
	return c.JSON(d)
}

func (h *EvaluationHandler) Approve(c *fiber.Ctx) error {
	adminID, ok := middleware.UserID(c)
	if !ok {
		return apperrors.ErrUnauthorized
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "geçersiz id")
	}
	d, err := h.svc.Approve(adminID, id)
	if err != nil {
		return err
	}
	return c.JSON(d)
}

func (h *EvaluationHandler) Revoke(c *fiber.Ctx) error {
	adminID, ok := middleware.UserID(c)
	if !ok {
		return apperrors.ErrUnauthorized
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "geçersiz id")
	}
	d, err := h.svc.Revoke(adminID, id)
	if err != nil {
		return err
	}
	return c.JSON(d)
}

func (h *EvaluationHandler) Versions(c *fiber.Ctx) error {
	adminID, ok := middleware.UserID(c)
	if !ok {
		return apperrors.ErrUnauthorized
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "geçersiz id")
	}
	role, _ := middleware.Role(c)
	if role != domain.RoleAdmin {
		return apperrors.ErrForbidden
	}
	items, err := h.svc.Versions(adminID, id)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *EvaluationHandler) MyStudents(c *fiber.Ctx) error {
	uid, _ := middleware.UserID(c)
	students, err := h.coach.MyStudents(uid)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": students})
}
