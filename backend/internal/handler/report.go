package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/koc-luk/backend/internal/domain"
	"github.com/koc-luk/backend/internal/middleware"
	"github.com/koc-luk/backend/internal/service"
	apperrors "github.com/koc-luk/backend/pkg/errors"
)

type ReportHandler struct {
	svc *service.ReportService
}

func NewReportHandler(svc *service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

func (h *ReportHandler) Summary(c *fiber.Ctx) error {
	r, err := h.svc.Summary()
	if err != nil {
		return err
	}
	return c.JSON(r)
}

func (h *ReportHandler) CoachPerformance(c *fiber.Ctx) error {
	rows, err := h.svc.CoachPerformance()
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": rows})
}

func (h *ReportHandler) CityDistribution(c *fiber.Ctx) error {
	role := c.Query("role", "student")
	rows, err := h.svc.CityDistribution(role)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": rows})
}

func (h *ReportHandler) WeekStats(c *fiber.Ctx) error {
	rows, err := h.svc.WeekStats()
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": rows})
}

func (h *ReportHandler) StudentTrend(c *fiber.Ctx) error {
	studentIDStr := c.Query("studentId")
	if studentIDStr == "" {
		return apperrors.WithDetails(apperrors.ErrValidation, "studentId gerekli")
	}
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "studentId geçersiz")
	}
	role, _ := middleware.Role(c)
	onlyApproved := role == domain.RoleParent || role == domain.RoleStudent
	rows, err := h.svc.StudentTrend(studentID, onlyApproved)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": rows})
}
