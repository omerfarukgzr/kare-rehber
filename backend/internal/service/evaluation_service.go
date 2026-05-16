package service

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/koc-luk/backend/internal/domain"
	"github.com/koc-luk/backend/internal/dto"
	"github.com/koc-luk/backend/internal/repository"
	apperrors "github.com/koc-luk/backend/pkg/errors"
)

type EvaluationService struct {
	db          *gorm.DB
	evals       *repository.EvaluationRepository
	weeks       *repository.WeekRepository
	assignments *repository.AssignmentRepository
	users       *repository.UserRepository
}

func NewEvaluationService(db *gorm.DB, evals *repository.EvaluationRepository, weeks *repository.WeekRepository, assignments *repository.AssignmentRepository, users *repository.UserRepository) *EvaluationService {
	return &EvaluationService{db: db, evals: evals, weeks: weeks, assignments: assignments, users: users}
}

func (s *EvaluationService) toDTO(row repository.EvaluationRow, role domain.Role) dto.EvaluationDTO {
	d := dto.EvaluationDTO{
		ID:           row.ID.String(),
		AssignmentID: row.AssignmentID.String(),
		WeekID:       row.EvaluationWeekID.String(),
		WeekNo:       row.WeekNo,
		WeekLabel:    row.WeekLabel,
		WeekIsOpen:   row.WeekIsOpen,
		CoachID:      row.CoachID.String(),
		CoachName:    row.CoachName,
		StudentID:    row.StudentID.String(),
		StudentName:  row.StudentName,
		StudentCity:  row.StudentCity,
		CourseStatus: row.CourseStatus,
		HomeworkDone: row.HomeworkDone,
		Motivation:   row.Motivation,
		Behavior:     row.Behavior,
		GeneralNote:  row.GeneralNote,
		Status:       row.Status,
		SubmittedBy:  row.SubmittedBy.String(),
		SubmittedAt:  row.SubmittedAt,
		ApprovedAt:   row.ApprovedAt,
		LastEditedAt: row.LastEditedAt,
	}
	if row.ApprovedBy != nil {
		s := row.ApprovedBy.String()
		d.ApprovedBy = &s
	}
	if row.LastEditedBy != nil {
		s := row.LastEditedBy.String()
		d.LastEditedBy = &s
	}
	if role == domain.RoleAdmin {
		d.AdminOnlyNote = row.AdminOnlyNote
	}
	return d
}

// Coach: kendi öğrencisi için yeni değerlendirme oluşturur
func (s *EvaluationService) CreateForCoach(coachID uuid.UUID, req dto.CreateEvaluationRequest) (*dto.EvaluationDTO, error) {
	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		return nil, apperrors.WithDetails(apperrors.ErrValidation, "student id geçersiz")
	}
	weekID, err := uuid.Parse(req.WeekID)
	if err != nil {
		return nil, apperrors.WithDetails(apperrors.ErrValidation, "week id geçersiz")
	}

	week, err := s.weeks.FindByID(weekID)
	if err != nil {
		return nil, err
	}
	if week == nil {
		return nil, apperrors.WithDetails(apperrors.ErrValidation, "hafta bulunamadı")
	}
	if !week.IsOpen {
		return nil, apperrors.WithDetails(apperrors.ErrConflict, "Bu hafta değerlendirmeye kapalı")
	}

	assignment, err := s.assignments.FindActiveByStudent(studentID)
	if err != nil {
		return nil, err
	}
	if assignment == nil || assignment.CoachID != coachID {
		return nil, apperrors.ErrForbidden
	}

	existing, err := s.evals.FindByAssignmentAndWeek(assignment.ID, weekID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, apperrors.WithDetails(apperrors.ErrConflict, "Bu hafta için bu öğrenciye zaten değerlendirme girilmiş")
	}

	now := time.Now()
	e := &domain.Evaluation{
		ID:               uuid.New(),
		AssignmentID:     assignment.ID,
		EvaluationWeekID: weekID,
		CourseStatus:     req.CourseStatus,
		HomeworkDone:     req.HomeworkDone,
		Motivation:       req.Motivation,
		Behavior:         req.Behavior,
		GeneralNote:      req.GeneralNote,
		AdminOnlyNote:    req.AdminOnlyNote,
		Status:           domain.EvalPending,
		SubmittedBy:      coachID,
		SubmittedAt:      now,
	}
	if err := s.evals.Create(e); err != nil {
		return nil, err
	}
	return s.GetByID(e.ID, domain.RoleAdmin)
}

// Coach: kendi değerlendirmesini düzenleyebilir (sadece pending iken)
func (s *EvaluationService) UpdateForCoach(coachID, evalID uuid.UUID, req dto.UpdateEvaluationRequest) (*dto.EvaluationDTO, error) {
	e, err := s.evals.FindByID(evalID)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, apperrors.ErrNotFound
	}
	assignment, err := s.assignments.FindByID(e.AssignmentID)
	if err != nil {
		return nil, err
	}
	if assignment == nil || assignment.CoachID != coachID {
		return nil, apperrors.ErrForbidden
	}
	if e.Status != domain.EvalPending {
		return nil, apperrors.WithDetails(apperrors.ErrConflict, "Onaylanmış veya yönetici düzenli değerlendirmeyi koç düzenleyemez")
	}
	week, err := s.weeks.FindByID(e.EvaluationWeekID)
	if err != nil {
		return nil, err
	}
	if week == nil || !week.IsOpen {
		return nil, apperrors.WithDetails(apperrors.ErrConflict, "Hafta kapalı")
	}

	applyUpdates(e, req)
	now := time.Now()
	e.LastEditedBy = &coachID
	e.LastEditedAt = &now
	if err := s.evals.Update(e); err != nil {
		return nil, err
	}
	return s.GetByID(e.ID, domain.RoleAdmin)
}

// Admin: doğrudan düzenleyebilir; eski sürüm version tablosuna yazılır
func (s *EvaluationService) UpdateByAdmin(adminID, evalID uuid.UUID, req dto.UpdateEvaluationRequest) (*dto.EvaluationDTO, error) {
	e, err := s.evals.FindByID(evalID)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, apperrors.ErrNotFound
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		nextVer, err := s.evals.NextVersionNo(e.ID)
		if err != nil {
			return err
		}
		snap := snapshotEvaluation(e)
		v := &domain.EvaluationVersion{
			ID:           uuid.New(),
			EvaluationID: e.ID,
			VersionNo:    nextVer,
			Snapshot:     snap,
			EditedBy:     &adminID,
			EditedAt:     time.Now(),
			ChangeReason: req.ChangeReason,
		}
		if err := tx.Create(v).Error; err != nil {
			return err
		}

		applyUpdates(e, req)
		now := time.Now()
		e.LastEditedBy = &adminID
		e.LastEditedAt = &now
		if e.Status == domain.EvalPending {
			e.Status = domain.EvalEditedByAdmin
		}
		if err := tx.Save(e).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.GetByID(e.ID, domain.RoleAdmin)
}

// Admin: onaylar
func (s *EvaluationService) Approve(adminID, evalID uuid.UUID) (*dto.EvaluationDTO, error) {
	e, err := s.evals.FindByID(evalID)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, apperrors.ErrNotFound
	}
	if e.Status == domain.EvalApproved {
		return s.GetByID(e.ID, domain.RoleAdmin)
	}
	now := time.Now()
	e.Status = domain.EvalApproved
	e.ApprovedBy = &adminID
	e.ApprovedAt = &now
	if err := s.evals.Update(e); err != nil {
		return nil, err
	}
	return s.GetByID(e.ID, domain.RoleAdmin)
}

// Admin: tekrar değerlendirilsin diye pending'e çekme (revoke)
func (s *EvaluationService) Revoke(adminID, evalID uuid.UUID) (*dto.EvaluationDTO, error) {
	e, err := s.evals.FindByID(evalID)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, apperrors.ErrNotFound
	}
	if e.Status == domain.EvalPending {
		return s.GetByID(e.ID, domain.RoleAdmin)
	}
	e.Status = domain.EvalPending
	e.ApprovedBy = nil
	e.ApprovedAt = nil
	now := time.Now()
	e.LastEditedBy = &adminID
	e.LastEditedAt = &now
	if err := s.evals.Update(e); err != nil {
		return nil, err
	}
	return s.GetByID(e.ID, domain.RoleAdmin)
}

func (s *EvaluationService) GetByID(id uuid.UUID, viewerRole domain.Role) (*dto.EvaluationDTO, error) {
	row, err := s.evals.FindRowByID(id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, apperrors.ErrNotFound
	}
	d := s.toDTO(*row, viewerRole)
	return &d, nil
}

// List with role-based filtering
func (s *EvaluationService) List(viewerID uuid.UUID, viewerRole domain.Role, q ListQuery) ([]dto.EvaluationDTO, int64, error) {
	p := repository.ListEvaluationParams{
		Search: q.Search,
		Offset: q.Offset, Limit: q.Limit,
	}
	if q.WeekID != "" {
		id, err := uuid.Parse(q.WeekID)
		if err != nil {
			return nil, 0, apperrors.WithDetails(apperrors.ErrValidation, "weekId geçersiz")
		}
		p.WeekID = &id
	}
	if q.CoachID != "" {
		id, err := uuid.Parse(q.CoachID)
		if err != nil {
			return nil, 0, apperrors.WithDetails(apperrors.ErrValidation, "coachId geçersiz")
		}
		p.CoachID = &id
	}
	if q.StudentID != "" {
		id, err := uuid.Parse(q.StudentID)
		if err != nil {
			return nil, 0, apperrors.WithDetails(apperrors.ErrValidation, "studentId geçersiz")
		}
		p.StudentID = &id
	}
	if q.Status != "" {
		st := domain.EvaluationStatus(q.Status)
		p.Status = &st
	}

	switch viewerRole {
	case domain.RoleAdmin:
	case domain.RoleCoach:
		p.CoachID = &viewerID
	case domain.RoleCoordinator:
		p.CoordinatorID = &viewerID
	case domain.RoleParent:
		p.ParentID = &viewerID
		p.OnlyApproved = true
	default:
		return nil, 0, apperrors.ErrForbidden
	}

	rows, total, err := s.evals.List(p)
	if err != nil {
		return nil, 0, err
	}
	items := make([]dto.EvaluationDTO, 0, len(rows))
	for _, r := range rows {
		items = append(items, s.toDTO(r, viewerRole))
	}
	return items, total, nil
}

type ListQuery struct {
	WeekID    string
	CoachID   string
	StudentID string
	Status    string
	Search    string
	Offset    int
	Limit     int
}

func (s *EvaluationService) Versions(adminID uuid.UUID, evalID uuid.UUID) ([]dto.EvaluationVersionDTO, error) {
	_ = adminID
	vs, err := s.evals.ListVersions(evalID)
	if err != nil {
		return nil, err
	}
	items := make([]dto.EvaluationVersionDTO, 0, len(vs))
	for i := range vs {
		items = append(items, dto.ToEvaluationVersionDTO(&vs[i]))
	}
	return items, nil
}

// Internal helpers

func applyUpdates(e *domain.Evaluation, req dto.UpdateEvaluationRequest) {
	if req.CourseStatus != nil {
		e.CourseStatus = req.CourseStatus
	}
	if req.HomeworkDone != nil {
		e.HomeworkDone = req.HomeworkDone
	}
	if req.Motivation != nil {
		e.Motivation = req.Motivation
	}
	if req.Behavior != nil {
		e.Behavior = req.Behavior
	}
	if req.GeneralNote != nil {
		e.GeneralNote = req.GeneralNote
	}
	if req.AdminOnlyNote != nil {
		e.AdminOnlyNote = req.AdminOnlyNote
	}
}

func snapshotEvaluation(e *domain.Evaluation) domain.JSONB {
	m := domain.JSONB{
		"courseStatus":  e.CourseStatus,
		"homeworkDone":  e.HomeworkDone,
		"motivation":    e.Motivation,
		"behavior":      e.Behavior,
		"generalNote":   e.GeneralNote,
		"adminOnlyNote": e.AdminOnlyNote,
		"status":        e.Status,
	}
	if e.LastEditedAt != nil {
		m["lastEditedAt"] = *e.LastEditedAt
	}
	return m
}
