package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/koc-luk/backend/internal/auth"
	"github.com/koc-luk/backend/internal/config"
	"github.com/koc-luk/backend/internal/domain"
	"github.com/koc-luk/backend/internal/dto"
	"github.com/koc-luk/backend/internal/repository"
	"github.com/koc-luk/backend/internal/sms"
	apperrors "github.com/koc-luk/backend/pkg/errors"
)

type UserService struct {
	db    *gorm.DB
	users *repository.UserRepository
	sms   sms.Provider
	cfg   *config.Config
}

func NewUserService(db *gorm.DB, users *repository.UserRepository, smsP sms.Provider, cfg *config.Config) *UserService {
	return &UserService{db: db, users: users, sms: smsP, cfg: cfg}
}

func (s *UserService) List(roleStr string, isActiveStr, city, search string, offset, limit int) ([]domain.User, int64, error) {
	var role *domain.Role
	if roleStr != "" {
		r := domain.Role(roleStr)
		role = &r
	}
	var isActive *bool
	if isActiveStr != "" {
		v := isActiveStr == "true"
		isActive = &v
	}
	var cityPtr *string
	if city != "" {
		cityPtr = &city
	}
	return s.users.ListByFilter(role, isActive, cityPtr, search, offset, limit)
}

func (s *UserService) Get(id uuid.UUID) (*domain.User, error) {
	u, err := s.users.FindByID(id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, apperrors.ErrNotFound
	}
	return u, nil
}

func (s *UserService) Create(ctx context.Context, actorID uuid.UUID, req dto.CreateUserRequest) (*domain.User, *string, error) {
	if !req.Role.Valid() {
		return nil, nil, apperrors.WithDetails(apperrors.ErrValidation, "geçersiz rol")
	}

	existing, err := s.users.FindByPhoneOrEmail(req.Phone)
	if err != nil {
		return nil, nil, err
	}
	if existing != nil {
		return nil, nil, apperrors.WithDetails(apperrors.ErrConflict, "Bu telefon numarasıyla bir kullanıcı zaten kayıtlı")
	}

	plain := req.Password
	autogen := false
	if plain == "" {
		p, err := auth.GenerateRandomPassword(10)
		if err != nil {
			return nil, nil, err
		}
		plain = p
		autogen = true
	}
	hash, err := auth.HashPassword(plain, s.cfg.BcryptCost)
	if err != nil {
		return nil, nil, err
	}

	u := &domain.User{
		ID:           uuid.New(),
		Role:         req.Role,
		FullName:     req.FullName,
		Phone:        req.Phone,
		PasswordHash: hash,
		IsActive:     true,
	}
	if req.Email != "" {
		e := req.Email
		u.Email = &e
	}
	if req.City != "" {
		c := req.City
		u.City = &c
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(u).Error; err != nil {
			return err
		}
		switch req.Role {
		case domain.RoleStudent:
			if err := tx.Create(&domain.Student{UserID: u.ID}).Error; err != nil {
				return err
			}
		case domain.RoleCoach:
			if err := tx.Create(&domain.Coach{UserID: u.ID, Status: domain.CoachApproved}).Error; err != nil {
				return err
			}
		case domain.RoleParent:
			if err := tx.Create(&domain.Parent{UserID: u.ID}).Error; err != nil {
				return err
			}
		case domain.RoleCoordinator:
			if err := tx.Create(&domain.Coordinator{UserID: u.ID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	if autogen {
		body := fmt.Sprintf("KARE-REHBER sistemine hoş geldiniz. Telefonunuz ve şu şifre ile giriş yapabilirsiniz: %s", plain)
		tplKey := "user_created"
		_, _ = s.sms.Send(ctx, sms.SendInput{
			ToPhone:     u.Phone,
			ToUserID:    &u.ID,
			Body:        body,
			SentBy:      &actorID,
			TemplateKey: &tplKey,
		})
		return u, &plain, nil
	}
	return u, nil, nil
}

func (s *UserService) Update(id uuid.UUID, req dto.UpdateUserRequest) (*domain.User, error) {
	u, err := s.users.FindByID(id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, apperrors.ErrNotFound
	}
	if req.FullName != nil && *req.FullName != "" {
		u.FullName = *req.FullName
	}
	if req.Email != nil {
		if *req.Email == "" {
			u.Email = nil
		} else {
			e := *req.Email
			u.Email = &e
		}
	}
	if req.City != nil {
		if *req.City == "" {
			u.City = nil
		} else {
			c := *req.City
			u.City = &c
		}
	}
	if req.IsActive != nil {
		u.IsActive = *req.IsActive
	}
	if err := s.users.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) ResetPassword(ctx context.Context, actorID, id uuid.UUID) (string, error) {
	u, err := s.users.FindByID(id)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", apperrors.ErrNotFound
	}
	plain, err := auth.GenerateRandomPassword(10)
	if err != nil {
		return "", err
	}
	hash, err := auth.HashPassword(plain, s.cfg.BcryptCost)
	if err != nil {
		return "", err
	}
	u.PasswordHash = hash
	if err := s.users.Update(u); err != nil {
		return "", err
	}
	body := fmt.Sprintf("KARE-REHBER şifreniz sıfırlandı. Yeni şifreniz: %s", plain)
	tplKey := "password_reset"
	_, _ = s.sms.Send(ctx, sms.SendInput{
		ToPhone:     u.Phone,
		ToUserID:    &u.ID,
		Body:        body,
		SentBy:      &actorID,
		TemplateKey: &tplKey,
	})
	return plain, nil
}
