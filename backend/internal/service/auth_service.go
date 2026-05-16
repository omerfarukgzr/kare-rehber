package service

import (
	"strings"

	"github.com/koc-luk/backend/internal/auth"
	"github.com/koc-luk/backend/internal/dto"
	"github.com/koc-luk/backend/internal/repository"
	apperrors "github.com/koc-luk/backend/pkg/errors"
)

type AuthService struct {
	users *repository.UserRepository
	jwt   *auth.Manager
}

func NewAuthService(users *repository.UserRepository, jwt *auth.Manager) *AuthService {
	return &AuthService{users: users, jwt: jwt}
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	identifier := strings.TrimSpace(req.Identifier)
	u, err := s.users.FindByPhoneOrEmail(identifier)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, apperrors.ErrInvalidCredential
	}
	if !u.IsActive {
		return nil, apperrors.ErrInvalidCredential
	}
	if err := auth.ComparePassword(u.PasswordHash, req.Password); err != nil {
		return nil, apperrors.ErrInvalidCredential
	}

	token, _, err := s.jwt.Generate(u.ID, u.Role)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User:  dto.ToUserDTO(u),
	}, nil
}
