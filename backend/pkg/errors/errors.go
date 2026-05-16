package errors

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type AppError struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	HTTPCode int    `json:"-"`
	Details  any    `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(httpCode int, code, message string) *AppError {
	return &AppError{Code: code, Message: message, HTTPCode: httpCode}
}

func WithDetails(err *AppError, details any) *AppError {
	return &AppError{Code: err.Code, Message: err.Message, HTTPCode: err.HTTPCode, Details: details}
}

var (
	ErrUnauthorized      = New(http.StatusUnauthorized, "unauthorized", "Yetkisiz erişim")
	ErrForbidden         = New(http.StatusForbidden, "forbidden", "Bu işlem için yetkiniz yok")
	ErrNotFound          = New(http.StatusNotFound, "not_found", "Kayıt bulunamadı")
	ErrConflict          = New(http.StatusConflict, "conflict", "Çakışan kayıt")
	ErrValidation        = New(http.StatusBadRequest, "validation_error", "Geçersiz veri")
	ErrInternal          = New(http.StatusInternalServerError, "internal_error", "Sunucu hatası")
	ErrInvalidCredential = New(http.StatusUnauthorized, "invalid_credentials", "Telefon/e-posta veya şifre hatalı")
)

func Handler(c *fiber.Ctx, err error) error {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return c.Status(appErr.HTTPCode).JSON(fiber.Map{"error": appErr})
	}

	if fe, ok := err.(*fiber.Error); ok {
		return c.Status(fe.Code).JSON(fiber.Map{"error": fiber.Map{"code": "http_error", "message": fe.Message}})
	}

	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"error": fiber.Map{"code": "internal_error", "message": "Sunucu hatası"},
	})
}
