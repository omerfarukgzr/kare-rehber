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

type AssignmentHandler struct {
	svc      *service.AssignmentService
	validate *validator.Validate
}

func NewAssignmentHandler(svc *service.AssignmentService, v *validator.Validate) *AssignmentHandler {
	return &AssignmentHandler{svc: svc, validate: v}
}

func (h *AssignmentHandler) List(c *fiber.Ctx) error {
	p := pagination.From(c)
	rows, total, err := h.svc.List(
		c.Query("coachId"),
		c.Query("studentId"),
		c.Query("city"),
		c.Query("search"),
		p.Offset(),
		p.Limit(),
	)
	if err != nil {
		return err
	}
	items := make([]dto.AssignmentDTO, 0, len(rows))
	for _, r := range rows {
		items = append(items, dto.AssignmentDTO{
			ID:           r.ID.String(),
			CoachID:      r.CoachID.String(),
			CoachName:    r.CoachName,
			CoachPhone:   r.CoachPhone,
			StudentID:    r.StudentID.String(),
			StudentName:  r.StudentName,
			StudentPhone: r.StudentPhone,
			StudentCity:  r.StudentCity,
			StartedAt:    r.StartedAt,
			EndedAt:      r.EndedAt,
			IsActive:     r.IsActive,
		})
	}
	return c.JSON(pagination.NewPage(items, total, p))
}

func (h *AssignmentHandler) Assign(c *fiber.Ctx) error {
	var req dto.AssignRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	by, ok := middleware.UserID(c)
	if !ok {
		return apperrors.ErrUnauthorized
	}
	res, err := h.svc.Assign(req, by)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (h *AssignmentHandler) Unassign(c *fiber.Ctx) error {
	var req dto.UnassignRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	if err := h.svc.Unassign(req.StudentID); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"ok": true})
}

func (h *AssignmentHandler) StudentsWithoutCoach(c *fiber.Ctx) error {
	users, err := h.svc.StudentsWithoutCoach(c.Query("city"), c.Query("search"))
	if err != nil {
		return err
	}
	items := make([]dto.UserDTO, 0, len(users))
	for i := range users {
		items = append(items, dto.ToUserDTO(&users[i]))
	}
	return c.JSON(fiber.Map{"items": items})
}

func (h *AssignmentHandler) SetCoordinator(c *fiber.Ctx) error {
	var req dto.SetCoordinatorRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	updated, err := h.svc.SetCoordinator(req)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"updated": updated})
}

func (h *AssignmentHandler) SetParent(c *fiber.Ctx) error {
	var req dto.SetParentRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, err.Error())
	}
	if err := h.svc.SetParent(req); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"ok": true})
}
