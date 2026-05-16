package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/koc-luk/backend/internal/domain"
	"github.com/koc-luk/backend/internal/dto"
	"github.com/koc-luk/backend/internal/repository"
	"github.com/koc-luk/backend/internal/sms"
	apperrors "github.com/koc-luk/backend/pkg/errors"
)

type SmsService struct {
	provider sms.Provider
	users    *repository.UserRepository
	evals    *repository.EvaluationRepository
	smsRepo  *repository.SmsRepository
}

func NewSmsService(p sms.Provider, users *repository.UserRepository, evals *repository.EvaluationRepository, smsRepo *repository.SmsRepository) *SmsService {
	return &SmsService{provider: p, users: users, evals: evals, smsRepo: smsRepo}
}

func (s *SmsService) Templates() []sms.Template { return sms.Templates }

func (s *SmsService) Send(ctx context.Context, actor uuid.UUID, req dto.SendSmsRequest) (*dto.SendSmsResponse, error) {
	if len(req.UserIDs) == 0 && len(req.Phones) == 0 {
		return nil, apperrors.WithDetails(apperrors.ErrValidation, "userIds veya phones gerekli")
	}
	resp := &dto.SendSmsResponse{}

	for _, uidStr := range req.UserIDs {
		uid, err := uuid.Parse(uidStr)
		if err != nil {
			resp.Failed++
			resp.Errors = append(resp.Errors, uidStr+": id geçersiz")
			continue
		}
		u, err := s.users.FindByID(uid)
		if err != nil {
			resp.Failed++
			resp.Errors = append(resp.Errors, uidStr+": "+err.Error())
			continue
		}
		if u == nil {
			resp.Failed++
			resp.Errors = append(resp.Errors, uidStr+": kullanıcı bulunamadı")
			continue
		}
		var tplKey *string
		if req.TemplateKey != "" {
			t := req.TemplateKey
			tplKey = &t
		}
		if _, err := s.provider.Send(ctx, sms.SendInput{
			ToPhone:     u.Phone,
			ToUserID:    &u.ID,
			Body:        req.Body,
			SentBy:      &actor,
			TemplateKey: tplKey,
		}); err != nil {
			resp.Failed++
			resp.Errors = append(resp.Errors, uidStr+": "+err.Error())
		} else {
			resp.Sent++
		}
	}

	for _, phone := range req.Phones {
		if phone == "" {
			continue
		}
		var tplKey *string
		if req.TemplateKey != "" {
			t := req.TemplateKey
			tplKey = &t
		}
		if _, err := s.provider.Send(ctx, sms.SendInput{
			ToPhone:     phone,
			Body:        req.Body,
			SentBy:      &actor,
			TemplateKey: tplKey,
		}); err != nil {
			resp.Failed++
			resp.Errors = append(resp.Errors, phone+": "+err.Error())
		} else {
			resp.Sent++
		}
	}

	return resp, nil
}

func (s *SmsService) Logs(search string, offset, limit int) ([]dto.SmsLogDTO, int64, error) {
	rows, total, err := s.smsRepo.List(search, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	items := make([]dto.SmsLogDTO, 0, len(rows))
	for _, r := range rows {
		d := dto.SmsLogDTO{
			ID:           r.ID.String(),
			ToPhone:      r.ToPhone,
			Body:         r.Body,
			Status:       string(r.Status),
			ProviderName: r.ProviderName,
			TemplateKey:  r.TemplateKey,
			SentAt:       r.SentAt,
		}
		if r.ToUserID != nil {
			s := r.ToUserID.String()
			d.ToUserID = &s
		}
		if r.SentByUserID != nil {
			s := r.SentByUserID.String()
			d.SentByUserID = &s
		}
		items = append(items, d)
	}
	return items, total, nil
}

func (s *SmsService) MissingCoachesForWeek(weekIDStr string) ([]dto.MissingCoachRow, error) {
	id, err := uuid.Parse(weekIDStr)
	if err != nil {
		return nil, apperrors.WithDetails(apperrors.ErrValidation, "weekId geçersiz")
	}
	rows, err := s.evals.MissingForWeek(id)
	if err != nil {
		return nil, err
	}
	items := make([]dto.MissingCoachRow, 0, len(rows))
	for _, r := range rows {
		items = append(items, dto.MissingCoachRow{
			CoachID:     r.CoachID.String(),
			CoachName:   r.CoachName,
			CoachPhone:  r.CoachPhone,
			StudentID:   r.StudentID.String(),
			StudentName: r.StudentName,
		})
	}
	return items, nil
}

// helper: ensure we use the ApprovalStatus enum somewhere; keep code compiling
var _ = domain.EvalApproved
