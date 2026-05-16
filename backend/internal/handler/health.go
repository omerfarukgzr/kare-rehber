package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Get(c *fiber.Ctx) error {
	dbStatus := "ok"
	if sqlDB, err := h.db.DB(); err != nil {
		dbStatus = "error"
	} else if err := sqlDB.Ping(); err != nil {
		dbStatus = "error"
	}

	return c.JSON(fiber.Map{
		"status":    "ok",
		"db":        dbStatus,
		"timestamp": time.Now().UTC(),
	})
}
