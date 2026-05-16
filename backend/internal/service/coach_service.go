package service

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/koc-luk/backend/internal/domain"
	"github.com/koc-luk/backend/internal/dto"
)

type CoachService struct {
	db *gorm.DB
}

func NewCoachService(db *gorm.DB) *CoachService {
	return &CoachService{db: db}
}

type CoachStudentRow struct {
	StudentID    uuid.UUID `json:"studentId" gorm:"column:student_id"`
	FullName     string    `json:"fullName" gorm:"column:full_name"`
	Phone        string    `json:"phone" gorm:"column:phone"`
	City         *string   `json:"city" gorm:"column:city"`
	AssignmentID uuid.UUID `json:"assignmentId" gorm:"column:assignment_id"`
}

func (s *CoachService) MyStudents(coachID uuid.UUID) ([]dto.UserDTO, error) {
	var users []domain.User
	err := s.db.Table("users AS u").
		Joins("JOIN coach_student_assignments a ON a.student_id = u.id AND a.is_active = true").
		Where("a.coach_id = ?", coachID).
		Order("u.full_name").
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	items := make([]dto.UserDTO, 0, len(users))
	for i := range users {
		items = append(items, dto.ToUserDTO(&users[i]))
	}
	return items, nil
}
