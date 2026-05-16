package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/koc-luk/backend/internal/domain"
)

type WeekRepository struct {
	db *gorm.DB
}

func NewWeekRepository(db *gorm.DB) *WeekRepository {
	return &WeekRepository{db: db}
}

func (r *WeekRepository) DB() *gorm.DB { return r.db }

func (r *WeekRepository) FindByID(id uuid.UUID) (*domain.EvaluationWeek, error) {
	var w domain.EvaluationWeek
	err := r.db.First(&w, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WeekRepository) FindByWeekNo(no int) (*domain.EvaluationWeek, error) {
	var w domain.EvaluationWeek
	err := r.db.First(&w, "week_no = ?", no).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WeekRepository) ListAll() ([]domain.EvaluationWeek, error) {
	var ws []domain.EvaluationWeek
	if err := r.db.Order("week_no").Find(&ws).Error; err != nil {
		return nil, err
	}
	return ws, nil
}

func (r *WeekRepository) ListOpen() ([]domain.EvaluationWeek, error) {
	var ws []domain.EvaluationWeek
	if err := r.db.Where("is_open = true").Order("week_no").Find(&ws).Error; err != nil {
		return nil, err
	}
	return ws, nil
}

func (r *WeekRepository) Create(w *domain.EvaluationWeek) error {
	return r.db.Create(w).Error
}

func (r *WeekRepository) Update(w *domain.EvaluationWeek) error {
	return r.db.Save(w).Error
}

// CurrentAndPreviousWeeks: returns the week containing today, and the one immediately before
func (r *WeekRepository) CurrentAndPreviousWeeks(now time.Time) ([]domain.EvaluationWeek, error) {
	var current domain.EvaluationWeek
	err := r.db.Where("? BETWEEN start_date AND end_date", now).First(&current).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Possibly outside academic year window — pick most recent past week as "current"
		err2 := r.db.Where("end_date <= ?", now).Order("end_date DESC").First(&current).Error
		if errors.Is(err2, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		if err2 != nil {
			return nil, err2
		}
	} else if err != nil {
		return nil, err
	}

	weeks := []domain.EvaluationWeek{current}
	if current.WeekNo > 1 {
		var prev domain.EvaluationWeek
		err := r.db.Where("week_no = ?", current.WeekNo-1).First(&prev).Error
		if err == nil {
			weeks = append(weeks, prev)
		}
	}
	return weeks, nil
}

// SyncOpenStates: aktif hafta + 1 önceki açık, diğerleri kapalı.
// Manuel açılmış (opened_by_user_id != NULL ve closed_at == NULL) haftaları açık tutar.
func (r *WeekRepository) SyncOpenStates(now time.Time) error {
	cur, err := r.CurrentAndPreviousWeeks(now)
	if err != nil {
		return err
	}
	openIDs := make([]uuid.UUID, 0, len(cur))
	for _, w := range cur {
		openIDs = append(openIDs, w.ID)
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	// Önce hepsini kapat (manuel açılmış olanlar hariç)
	closeQ := tx.Model(&domain.EvaluationWeek{}).Where("is_open = true")
	if len(openIDs) > 0 {
		closeQ = closeQ.Where("id NOT IN ?", openIDs)
	}
	closeQ = closeQ.Where("opened_by_user_id IS NULL OR closed_at IS NOT NULL")
	if err := closeQ.Updates(map[string]any{"is_open": false, "closed_at": now}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(openIDs) > 0 {
		if err := tx.Model(&domain.EvaluationWeek{}).
			Where("id IN ?", openIDs).
			Updates(map[string]any{"is_open": true}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
