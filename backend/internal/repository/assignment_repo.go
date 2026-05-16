package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/koc-luk/backend/internal/domain"
)

type AssignmentRepository struct {
	db *gorm.DB
}

func NewAssignmentRepository(db *gorm.DB) *AssignmentRepository {
	return &AssignmentRepository{db: db}
}

func (r *AssignmentRepository) DB() *gorm.DB { return r.db }

func (r *AssignmentRepository) FindActiveByStudent(studentID uuid.UUID) (*domain.CoachStudentAssignment, error) {
	var a domain.CoachStudentAssignment
	err := r.db.Where("student_id = ? AND is_active = true", studentID).First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AssignmentRepository) FindByID(id uuid.UUID) (*domain.CoachStudentAssignment, error) {
	var a domain.CoachStudentAssignment
	err := r.db.First(&a, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AssignmentRepository) Assign(coachID, studentID uuid.UUID, by *uuid.UUID) (*domain.CoachStudentAssignment, error) {
	now := time.Now()
	var assignment *domain.CoachStudentAssignment
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.CoachStudentAssignment{}).
			Where("student_id = ? AND is_active = true", studentID).
			Updates(map[string]any{"is_active": false, "ended_at": &now}).Error; err != nil {
			return err
		}
		a := &domain.CoachStudentAssignment{
			ID:         uuid.New(),
			CoachID:    coachID,
			StudentID:  studentID,
			StartedAt:  now,
			IsActive:   true,
			AssignedBy: by,
		}
		if err := tx.Create(a).Error; err != nil {
			return err
		}
		assignment = a
		return nil
	})
	return assignment, err
}

func (r *AssignmentRepository) Unassign(studentID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&domain.CoachStudentAssignment{}).
		Where("student_id = ? AND is_active = true", studentID).
		Updates(map[string]any{"is_active": false, "ended_at": &now}).Error
}

type AssignmentRow struct {
	domain.CoachStudentAssignment
	CoachName   string `json:"coachName" gorm:"column:coach_name"`
	CoachPhone  string `json:"coachPhone" gorm:"column:coach_phone"`
	StudentName string `json:"studentName" gorm:"column:student_name"`
	StudentPhone string `json:"studentPhone" gorm:"column:student_phone"`
	StudentCity *string `json:"studentCity" gorm:"column:student_city"`
}

func (r *AssignmentRepository) ListActive(coachID, studentID *uuid.UUID, city *string, search string, offset, limit int) ([]AssignmentRow, int64, error) {
	q := r.db.Table("coach_student_assignments AS a").
		Select(`a.*,
			c.full_name AS coach_name, c.phone AS coach_phone,
			s.full_name AS student_name, s.phone AS student_phone, s.city AS student_city`).
		Joins("JOIN users c ON c.id = a.coach_id").
		Joins("JOIN users s ON s.id = a.student_id").
		Where("a.is_active = true")

	if coachID != nil {
		q = q.Where("a.coach_id = ?", *coachID)
	}
	if studentID != nil {
		q = q.Where("a.student_id = ?", *studentID)
	}
	if city != nil && *city != "" {
		q = q.Where("s.city = ?", *city)
	}
	if search != "" {
		s := "%" + search + "%"
		q = q.Where("c.full_name ILIKE ? OR s.full_name ILIKE ?", s, s)
	}

	var total int64
	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []AssignmentRow
	if err := q.Order("a.created_at DESC").Offset(offset).Limit(limit).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// StudentsWithoutCoach: aktif eşleştirmesi olmayan öğrenciler (is_active=true users with role=student)
func (r *AssignmentRepository) StudentsWithoutCoach(city *string, search string) ([]domain.User, error) {
	q := r.db.Table("users AS u").
		Where("u.role = 'student' AND u.is_active = true").
		Where("NOT EXISTS (SELECT 1 FROM coach_student_assignments a WHERE a.student_id = u.id AND a.is_active = true)")
	if city != nil && *city != "" {
		q = q.Where("u.city = ?", *city)
	}
	if search != "" {
		s := "%" + search + "%"
		q = q.Where("u.full_name ILIKE ? OR u.phone ILIKE ?", s, s)
	}
	var users []domain.User
	if err := q.Order("u.created_at DESC").Limit(500).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
