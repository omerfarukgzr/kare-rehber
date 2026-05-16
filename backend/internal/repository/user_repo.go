package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/koc-luk/backend/internal/domain"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) DB() *gorm.DB { return r.db }

func (r *UserRepository) FindByPhoneOrEmail(identifier string) (*domain.User, error) {
	var u domain.User
	err := r.db.Where("phone = ? OR email = ?", identifier, identifier).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	var u domain.User
	err := r.db.First(&u, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Create(u *domain.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepository) Update(u *domain.User) error {
	return r.db.Save(u).Error
}

func (r *UserRepository) ListByFilter(role *domain.Role, isActive *bool, city *string, search string, offset, limit int) ([]domain.User, int64, error) {
	q := r.db.Model(&domain.User{})
	if role != nil {
		q = q.Where("role = ?", *role)
	}
	if isActive != nil {
		q = q.Where("is_active = ?", *isActive)
	}
	if city != nil && *city != "" {
		q = q.Where("city = ?", *city)
	}
	if search != "" {
		s := "%" + search + "%"
		q = q.Where("full_name ILIKE ? OR phone ILIKE ? OR email ILIKE ?", s, s, s)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []domain.User
	if err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
