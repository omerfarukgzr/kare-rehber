package service

import (
	"context"
	"errors"
	"fmt"
	"time"

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

type RegistrationService struct {
	db    *gorm.DB
	regs  *repository.RegistrationRepository
	users *repository.UserRepository
	sms   sms.Provider
	cfg   *config.Config
}

func NewRegistrationService(db *gorm.DB, regs *repository.RegistrationRepository, users *repository.UserRepository, smsP sms.Provider, cfg *config.Config) *RegistrationService {
	return &RegistrationService{db: db, regs: regs, users: users, sms: smsP, cfg: cfg}
}

func (s *RegistrationService) ApplyStudent(req dto.StudentRegistrationRequest) (*domain.Registration, error) {
	payload := domain.JSONB{
		"school":         req.School,
		"grade":          req.Grade,
		"parentFullName": req.ParentFullName,
		"parentPhone":    req.ParentPhone,
		"notes":          req.Notes,
	}
	return s.applyCommon(domain.RegistrationKindStudent, req.FullName, req.Phone, req.Email, req.City, payload)
}

func (s *RegistrationService) ApplyCoach(req dto.CoachRegistrationRequest) (*domain.Registration, error) {
	payload := domain.JSONB{
		"bio":        req.Bio,
		"experience": req.Experience,
	}
	return s.applyCommon(domain.RegistrationKindCoach, req.FullName, req.Phone, req.Email, req.City, payload)
}

func (s *RegistrationService) applyCommon(kind domain.RegistrationKind, fullName, phone, email, city string, payload domain.JSONB) (*domain.Registration, error) {
	existing, err := s.users.FindByPhoneOrEmail(phone)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, apperrors.WithDetails(apperrors.ErrConflict, "Bu telefon numarasıyla bir kullanıcı zaten kayıtlı")
	}

	reg := &domain.Registration{
		ID:       uuid.New(),
		Kind:     kind,
		Status:   domain.RegistrationPending,
		Payload:  payload,
		FullName: fullName,
		Phone:    phone,
	}
	if email != "" {
		e := email
		reg.Email = &e
	}
	if city != "" {
		c := city
		reg.City = &c
	}
	if err := s.regs.Create(reg); err != nil {
		return nil, err
	}
	return reg, nil
}

func (s *RegistrationService) List(kindStr, statusStr, search string, offset, limit int) ([]domain.Registration, int64, error) {
	var kind *domain.RegistrationKind
	if kindStr != "" {
		k := domain.RegistrationKind(kindStr)
		kind = &k
	}
	var status *domain.RegistrationStatus
	if statusStr != "" {
		st := domain.RegistrationStatus(statusStr)
		status = &st
	}
	return s.regs.List(kind, status, search, offset, limit)
}

func (s *RegistrationService) Decide(ctx context.Context, regID uuid.UUID, actorID uuid.UUID, req dto.RegistrationDecisionRequest) (*dto.RegistrationDecisionResponse, error) {
	reg, err := s.regs.FindByID(regID)
	if err != nil {
		return nil, err
	}
	if reg == nil {
		return nil, apperrors.ErrNotFound
	}
	if reg.Status != domain.RegistrationPending {
		return nil, apperrors.WithDetails(apperrors.ErrConflict, "Bu başvuru zaten karara bağlanmış")
	}

	now := time.Now()
	resp := &dto.RegistrationDecisionResponse{}

	if !req.Approve {
		reg.Status = domain.RegistrationRejected
		reg.ReviewedBy = &actorID
		reg.ReviewedAt = &now
		if req.Note != "" {
			n := req.Note
			reg.ReviewNote = &n
		}
		if err := s.regs.Update(reg); err != nil {
			return nil, err
		}
		resp.Registration = dto.ToRegistrationDTO(reg)
		return resp, nil
	}

	plain, err := auth.GenerateRandomPassword(10)
	if err != nil {
		return nil, err
	}
	hash, err := auth.HashPassword(plain, s.cfg.BcryptCost)
	if err != nil {
		return nil, err
	}

	var role domain.Role
	switch reg.Kind {
	case domain.RegistrationKindStudent:
		role = domain.RoleStudent
	case domain.RegistrationKindCoach:
		role = domain.RoleCoach
	default:
		return nil, errors.New("unknown registration kind")
	}

	user := &domain.User{
		ID:           uuid.New(),
		Role:         role,
		FullName:     reg.FullName,
		Phone:        reg.Phone,
		Email:        reg.Email,
		PasswordHash: hash,
		IsActive:     true,
		City:         reg.City,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		switch reg.Kind {
		case domain.RegistrationKindStudent:
			payload := reg.Payload
			school := stringFromMap(payload, "school")
			grade := stringFromMap(payload, "grade")
			notes := stringFromMap(payload, "notes")
			st := &domain.Student{
				UserID: user.ID,
				School: school,
				Grade:  grade,
				Notes:  notes,
			}
			if err := tx.Create(st).Error; err != nil {
				return err
			}
		case domain.RegistrationKindCoach:
			payload := reg.Payload
			bio := stringFromMap(payload, "bio")
			experience := stringFromMap(payload, "experience")
			co := &domain.Coach{
				UserID:     user.ID,
				Status:     domain.CoachApproved,
				Bio:        bio,
				Experience: experience,
			}
			if err := tx.Create(co).Error; err != nil {
				return err
			}
		}

		reg.Status = domain.RegistrationApproved
		reg.ReviewedBy = &actorID
		reg.ReviewedAt = &now
		reg.CreatedUserID = &user.ID
		if req.Note != "" {
			n := req.Note
			reg.ReviewNote = &n
		}
		if err := tx.Save(reg).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	body := fmt.Sprintf("KARE-REHBER sistemine hoş geldiniz. Telefon numaranız ve şu şifre ile giriş yapabilirsiniz: %s", plain)
	tplKey := "registration_approved"
	smsRes, smsErr := s.sms.Send(ctx, sms.SendInput{
		ToPhone:     user.Phone,
		ToUserID:    &user.ID,
		Body:        body,
		SentBy:      &actorID,
		TemplateKey: &tplKey,
	})

	uid := user.ID.String()
	resp.Registration = dto.ToRegistrationDTO(reg)
	resp.UserID = &uid
	resp.GeneratedPwd = &plain
	if smsErr == nil && smsRes != nil {
		s := smsRes.LogID.String()
		resp.SmsLogID = &s
	}
	return resp, nil
}

func stringFromMap(m domain.JSONB, key string) *string {
	if m == nil {
		return nil
	}
	v, ok := m[key]
	if !ok {
		return nil
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return nil
	}
	return &s
}
