package dto

import (
	"time"

	"github.com/koc-luk/backend/internal/domain"
)

type StudentRegistrationRequest struct {
	FullName       string `json:"fullName" validate:"required"`
	Phone          string `json:"phone" validate:"required"`
	Email          string `json:"email"`
	City           string `json:"city"`
	School         string `json:"school"`
	Grade          string `json:"grade"`
	ParentFullName string `json:"parentFullName"`
	ParentPhone    string `json:"parentPhone"`
	Notes          string `json:"notes"`
}

type CoachRegistrationRequest struct {
	FullName   string `json:"fullName" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
	Email      string `json:"email"`
	City       string `json:"city"`
	Bio        string `json:"bio"`
	Experience string `json:"experience"`
}

type RegistrationDTO struct {
	ID         string                    `json:"id"`
	Kind       domain.RegistrationKind   `json:"kind"`
	Status     domain.RegistrationStatus `json:"status"`
	FullName   string                    `json:"fullName"`
	Phone      string                    `json:"phone"`
	Email      *string                   `json:"email"`
	City       *string                   `json:"city"`
	Payload    domain.JSONB              `json:"payload"`
	ReviewedAt *time.Time                `json:"reviewedAt"`
	ReviewNote *string                   `json:"reviewNote"`
	CreatedAt  time.Time                 `json:"createdAt"`
}

type RegistrationDecisionRequest struct {
	Approve bool   `json:"approve"`
	Note    string `json:"note"`
}

type RegistrationDecisionResponse struct {
	Registration RegistrationDTO `json:"registration"`
	UserID       *string         `json:"userId,omitempty"`
	GeneratedPwd *string         `json:"generatedPassword,omitempty"`
	SmsLogID     *string         `json:"smsLogId,omitempty"`
}

func ToRegistrationDTO(r *domain.Registration) RegistrationDTO {
	return RegistrationDTO{
		ID:         r.ID.String(),
		Kind:       r.Kind,
		Status:     r.Status,
		FullName:   r.FullName,
		Phone:      r.Phone,
		Email:      r.Email,
		City:       r.City,
		Payload:    r.Payload,
		ReviewedAt: r.ReviewedAt,
		ReviewNote: r.ReviewNote,
		CreatedAt:  r.CreatedAt,
	}
}
