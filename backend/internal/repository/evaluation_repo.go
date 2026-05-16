package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/koc-luk/backend/internal/domain"
)

type EvaluationRepository struct {
	db *gorm.DB
}

func NewEvaluationRepository(db *gorm.DB) *EvaluationRepository {
	return &EvaluationRepository{db: db}
}

func (r *EvaluationRepository) DB() *gorm.DB { return r.db }

func (r *EvaluationRepository) FindByID(id uuid.UUID) (*domain.Evaluation, error) {
	var e domain.Evaluation
	err := r.db.First(&e, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *EvaluationRepository) FindRowByID(id uuid.UUID) (*EvaluationRow, error) {
	var row EvaluationRow
	q := r.db.Table("evaluations AS e").
		Select(`e.*,
			w.week_no AS week_no, w.label AS week_label, w.is_open AS week_is_open,
			a.coach_id AS coach_id, c.full_name AS coach_name,
			a.student_id AS student_id, s.full_name AS student_name, s.city AS student_city`).
		Joins("JOIN evaluation_weeks w ON w.id = e.evaluation_week_id").
		Joins("JOIN coach_student_assignments a ON a.id = e.assignment_id").
		Joins("JOIN users c ON c.id = a.coach_id").
		Joins("JOIN users s ON s.id = a.student_id").
		Where("e.id = ?", id)
	err := q.Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == uuid.Nil {
		return nil, nil
	}
	return &row, nil
}

func (r *EvaluationRepository) FindByAssignmentAndWeek(assignmentID, weekID uuid.UUID) (*domain.Evaluation, error) {
	var e domain.Evaluation
	err := r.db.Where("assignment_id = ? AND evaluation_week_id = ?", assignmentID, weekID).First(&e).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *EvaluationRepository) Create(e *domain.Evaluation) error {
	return r.db.Create(e).Error
}

func (r *EvaluationRepository) Update(e *domain.Evaluation) error {
	return r.db.Save(e).Error
}

func (r *EvaluationRepository) NextVersionNo(evaluationID uuid.UUID) (int, error) {
	var max int
	row := r.db.Model(&domain.EvaluationVersion{}).
		Where("evaluation_id = ?", evaluationID).
		Select("COALESCE(MAX(version_no), 0) AS max")
	if err := row.Scan(&max).Error; err != nil {
		return 0, err
	}
	return max + 1, nil
}

func (r *EvaluationRepository) AppendVersion(v *domain.EvaluationVersion) error {
	return r.db.Create(v).Error
}

func (r *EvaluationRepository) ListVersions(evaluationID uuid.UUID) ([]domain.EvaluationVersion, error) {
	var vs []domain.EvaluationVersion
	if err := r.db.Where("evaluation_id = ?", evaluationID).Order("version_no DESC").Find(&vs).Error; err != nil {
		return nil, err
	}
	return vs, nil
}

type EvaluationRow struct {
	domain.Evaluation
	WeekNo       int    `json:"weekNo" gorm:"column:week_no"`
	WeekLabel    string `json:"weekLabel" gorm:"column:week_label"`
	WeekIsOpen   bool   `json:"weekIsOpen" gorm:"column:week_is_open"`
	CoachID      uuid.UUID `json:"coachId" gorm:"column:coach_id"`
	CoachName    string `json:"coachName" gorm:"column:coach_name"`
	StudentID    uuid.UUID `json:"studentId" gorm:"column:student_id"`
	StudentName  string `json:"studentName" gorm:"column:student_name"`
	StudentCity  *string `json:"studentCity" gorm:"column:student_city"`
}

type ListEvaluationParams struct {
	WeekID         *uuid.UUID
	CoachID        *uuid.UUID
	StudentID      *uuid.UUID
	CoordinatorID  *uuid.UUID // student.coordinator_id eşleşenler
	ParentID       *uuid.UUID // student.parent_id eşleşenler
	OnlyApproved   bool
	Status         *domain.EvaluationStatus
	Search         string
	Offset         int
	Limit          int
}

func (r *EvaluationRepository) List(p ListEvaluationParams) ([]EvaluationRow, int64, error) {
	q := r.db.Table("evaluations AS e").
		Select(`e.*,
			w.week_no AS week_no, w.label AS week_label, w.is_open AS week_is_open,
			a.coach_id AS coach_id, c.full_name AS coach_name,
			a.student_id AS student_id, s.full_name AS student_name, s.city AS student_city`).
		Joins("JOIN evaluation_weeks w ON w.id = e.evaluation_week_id").
		Joins("JOIN coach_student_assignments a ON a.id = e.assignment_id").
		Joins("JOIN users c ON c.id = a.coach_id").
		Joins("JOIN users s ON s.id = a.student_id")

	if p.CoordinatorID != nil {
		q = q.Joins("JOIN students st ON st.user_id = a.student_id").
			Where("st.coordinator_id = ?", *p.CoordinatorID)
	}
	if p.ParentID != nil {
		q = q.Joins("JOIN students stp ON stp.user_id = a.student_id").
			Where("stp.parent_id = ?", *p.ParentID)
	}

	if p.WeekID != nil {
		q = q.Where("e.evaluation_week_id = ?", *p.WeekID)
	}
	if p.CoachID != nil {
		q = q.Where("a.coach_id = ?", *p.CoachID)
	}
	if p.StudentID != nil {
		q = q.Where("a.student_id = ?", *p.StudentID)
	}
	if p.OnlyApproved {
		q = q.Where("e.status = ?", domain.EvalApproved)
	}
	if p.Status != nil {
		q = q.Where("e.status = ?", *p.Status)
	}
	if p.Search != "" {
		s := "%" + p.Search + "%"
		q = q.Where("c.full_name ILIKE ? OR s.full_name ILIKE ?", s, s)
	}

	var total int64
	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []EvaluationRow
	if err := q.Order("w.week_no DESC, e.created_at DESC").Offset(p.Offset).Limit(p.Limit).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// MissingForWeek: bu hafta için değerlendirme girilmemiş aktif eşleştirmelerin koçları ile öğrencileri
type MissingRow struct {
	AssignmentID uuid.UUID `gorm:"column:assignment_id" json:"assignmentId"`
	CoachID      uuid.UUID `gorm:"column:coach_id" json:"coachId"`
	CoachName    string    `gorm:"column:coach_name" json:"coachName"`
	CoachPhone   string    `gorm:"column:coach_phone" json:"coachPhone"`
	StudentID    uuid.UUID `gorm:"column:student_id" json:"studentId"`
	StudentName  string    `gorm:"column:student_name" json:"studentName"`
}

func (r *EvaluationRepository) MissingForWeek(weekID uuid.UUID) ([]MissingRow, error) {
	var rows []MissingRow
	q := r.db.Table("coach_student_assignments AS a").
		Select(`a.id AS assignment_id,
			a.coach_id AS coach_id, c.full_name AS coach_name, c.phone AS coach_phone,
			a.student_id AS student_id, s.full_name AS student_name`).
		Joins("JOIN users c ON c.id = a.coach_id AND c.is_active = true").
		Joins("JOIN users s ON s.id = a.student_id AND s.is_active = true").
		Where("a.is_active = true").
		Where("NOT EXISTS (SELECT 1 FROM evaluations e WHERE e.assignment_id = a.id AND e.evaluation_week_id = ?)", weekID).
		Order("c.full_name, s.full_name")
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
