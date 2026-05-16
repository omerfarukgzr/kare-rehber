package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/koc-luk/backend/internal/domain"
	"github.com/koc-luk/backend/internal/dto"
	"github.com/koc-luk/backend/internal/repository"
	apperrors "github.com/koc-luk/backend/pkg/errors"
)

type WeekService struct {
	db    *gorm.DB
	weeks *repository.WeekRepository
}

func NewWeekService(db *gorm.DB, weeks *repository.WeekRepository) *WeekService {
	return &WeekService{db: db, weeks: weeks}
}

const dateLayout = "2006-01-02"

func parseDate(s string) (time.Time, error) {
	return time.ParseInLocation(dateLayout, s, time.Local)
}

func (s *WeekService) ListAll() ([]domain.EvaluationWeek, error) {
	if err := s.weeks.SyncOpenStates(time.Now()); err != nil {
		return nil, err
	}
	return s.weeks.ListAll()
}

func (s *WeekService) ListOpen() ([]domain.EvaluationWeek, error) {
	if err := s.weeks.SyncOpenStates(time.Now()); err != nil {
		return nil, err
	}
	return s.weeks.ListOpen()
}

func (s *WeekService) Create(req dto.CreateWeekRequest) (*domain.EvaluationWeek, error) {
	start, err := parseDate(req.StartDate)
	if err != nil {
		return nil, apperrors.WithDetails(apperrors.ErrValidation, "başlangıç tarihi formatı YYYY-MM-DD olmalı")
	}
	end, err := parseDate(req.EndDate)
	if err != nil {
		return nil, apperrors.WithDetails(apperrors.ErrValidation, "bitiş tarihi formatı YYYY-MM-DD olmalı")
	}
	if end.Before(start) {
		return nil, apperrors.WithDetails(apperrors.ErrValidation, "bitiş başlangıçtan önce olamaz")
	}
	existing, _ := s.weeks.FindByWeekNo(req.WeekNo)
	if existing != nil {
		return nil, apperrors.WithDetails(apperrors.ErrConflict, fmt.Sprintf("Hafta %d zaten mevcut", req.WeekNo))
	}
	label := req.Label
	if label == "" {
		label = fmt.Sprintf("%d. Hafta (%s-%s)", req.WeekNo, start.Format("02 Jan"), end.Format("02 Jan"))
	}
	w := &domain.EvaluationWeek{
		ID:        uuid.New(),
		WeekNo:    req.WeekNo,
		Label:     label,
		StartDate: start,
		EndDate:   end,
		IsOpen:    false,
	}
	if err := s.weeks.Create(w); err != nil {
		return nil, err
	}
	_ = s.weeks.SyncOpenStates(time.Now())
	return w, nil
}

func (s *WeekService) Update(id uuid.UUID, req dto.UpdateWeekRequest) (*domain.EvaluationWeek, error) {
	w, err := s.weeks.FindByID(id)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, apperrors.ErrNotFound
	}
	if req.Label != nil {
		w.Label = *req.Label
	}
	if req.StartDate != nil {
		t, err := parseDate(*req.StartDate)
		if err != nil {
			return nil, apperrors.WithDetails(apperrors.ErrValidation, "başlangıç tarihi geçersiz")
		}
		w.StartDate = t
	}
	if req.EndDate != nil {
		t, err := parseDate(*req.EndDate)
		if err != nil {
			return nil, apperrors.WithDetails(apperrors.ErrValidation, "bitiş tarihi geçersiz")
		}
		w.EndDate = t
	}
	if err := s.weeks.Update(w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WeekService) Open(id uuid.UUID, by uuid.UUID) (*domain.EvaluationWeek, error) {
	w, err := s.weeks.FindByID(id)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, apperrors.ErrNotFound
	}
	now := time.Now()
	w.IsOpen = true
	w.OpenedByUserID = &by
	w.OpenedAt = &now
	w.ClosedAt = nil
	if err := s.weeks.Update(w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WeekService) Close(id uuid.UUID) (*domain.EvaluationWeek, error) {
	w, err := s.weeks.FindByID(id)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, apperrors.ErrNotFound
	}
	now := time.Now()
	w.IsOpen = false
	w.ClosedAt = &now
	if err := s.weeks.Update(w); err != nil {
		return nil, err
	}
	return w, nil
}

// Generate: belirtilen tarihten itibaren n hafta üretir (week_no'lar 1'den başlar veya devam eder)
func (s *WeekService) Generate(req dto.GenerateWeeksRequest) ([]domain.EvaluationWeek, error) {
	start, err := parseDate(req.StartDate)
	if err != nil {
		return nil, apperrors.WithDetails(apperrors.ErrValidation, "başlangıç tarihi geçersiz")
	}

	existing, err := s.weeks.ListAll()
	if err != nil {
		return nil, err
	}
	startNo := 1
	for _, w := range existing {
		if w.WeekNo >= startNo {
			startNo = w.WeekNo + 1
		}
	}

	created := make([]domain.EvaluationWeek, 0, req.WeekCount)
	err = s.db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < req.WeekCount; i++ {
			no := startNo + i
			ws := start.AddDate(0, 0, i*7)
			we := ws.AddDate(0, 0, 6)
			w := domain.EvaluationWeek{
				ID:        uuid.New(),
				WeekNo:    no,
				Label:     fmt.Sprintf("%d. Hafta (%s-%s)", no, ws.Format("02 Jan"), we.Format("02 Jan")),
				StartDate: ws,
				EndDate:   we,
				IsOpen:    false,
			}
			if err := tx.Create(&w).Error; err != nil {
				return err
			}
			created = append(created, w)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	_ = s.weeks.SyncOpenStates(time.Now())
	return created, nil
}
