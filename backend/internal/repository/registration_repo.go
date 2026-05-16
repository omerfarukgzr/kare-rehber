package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/koc-luk/backend/internal/domain"
)

type RegistrationRepository struct {
	db *gorm.DB
}

func NewRegistrationRepository(db *gorm.DB) *RegistrationRepository {
	return &RegistrationRepository{db: db}
}

func (r *RegistrationRepository) Create(reg *domain.Registration) error {
	return r.db.Create(reg).Error
}

func (r *RegistrationRepository) FindByID(id uuid.UUID) (*domain.Registration, error) {
	var reg domain.Registration
	err := r.db.First(&reg, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func (r *RegistrationRepository) Update(reg *domain.Registration) error {
	return r.db.Save(reg).Error
}

func (r *RegistrationRepository) List(kind *domain.RegistrationKind, status *domain.RegistrationStatus, search string, offset, limit int) ([]domain.Registration, int64, error) {
	q := r.db.Model(&domain.Registration{})
	if kind != nil {
		q = q.Where("kind = ?", *kind)
	}
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	if search != "" {
		s := "%" + search + "%"
		q = q.Where("full_name ILIKE ? OR phone ILIKE ? OR email ILIKE ?", s, s, s)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []domain.Registration
	if err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
