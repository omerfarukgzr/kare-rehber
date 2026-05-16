package dto

import (
	"github.com/koc-luk/backend/internal/domain"
)

type CreateUserRequest struct {
	Role     domain.Role `json:"role" validate:"required,oneof=admin coach student parent coordinator"`
	FullName string      `json:"fullName" validate:"required"`
	Phone    string      `json:"phone" validate:"required"`
	Email    string      `json:"email"`
	City     string      `json:"city"`
	Password string      `json:"password"`
}

type UpdateUserRequest struct {
	FullName *string `json:"fullName"`
	Email    *string `json:"email"`
	City     *string `json:"city"`
	IsActive *bool   `json:"isActive"`
}

type CreateUserResponse struct {
	User              UserDTO `json:"user"`
	GeneratedPassword *string `json:"generatedPassword,omitempty"`
}
