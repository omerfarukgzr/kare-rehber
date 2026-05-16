package repository

import (
	"gorm.io/gorm"

	"github.com/koc-luk/backend/internal/domain"
)

type SmsRepository struct {
	db *gorm.DB
}

func NewSmsRepository(db *gorm.DB) *SmsRepository {
	return &SmsRepository{db: db}
}

func (r *SmsRepository) List(search string, offset, limit int) ([]domain.SmsLog, int64, error) {
	q := r.db.Model(&domain.SmsLog{})
	if search != "" {
		s := "%" + search + "%"
		q = q.Where("to_phone ILIKE ? OR body ILIKE ?", s, s)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []domain.SmsLog
	if err := q.Order("sent_at DESC").Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}
