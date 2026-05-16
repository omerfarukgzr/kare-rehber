package dto

import "github.com/koc-luk/backend/internal/domain"

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required,min=4"`
}

type UserDTO struct {
	ID        string      `json:"id"`
	Role      domain.Role `json:"role"`
	FullName  string      `json:"fullName"`
	Phone     string      `json:"phone"`
	Email     *string     `json:"email"`
	City      *string     `json:"city"`
	IsActive  bool        `json:"isActive"`
	CreatedAt string      `json:"createdAt"`
}

type LoginResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

func ToUserDTO(u *domain.User) UserDTO {
	return UserDTO{
		ID:        u.ID.String(),
		Role:      u.Role,
		FullName:  u.FullName,
		Phone:     u.Phone,
		Email:     u.Email,
		City:      u.City,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
