package service

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/koc-luk/backend/internal/domain"
	"github.com/koc-luk/backend/internal/dto"
	"github.com/koc-luk/backend/internal/repository"
	apperrors "github.com/koc-luk/backend/pkg/errors"
)

type AssignmentService struct {
	db          *gorm.DB
	assignments *repository.AssignmentRepository
	users       *repository.UserRepository
}

func NewAssignmentService(db *gorm.DB, assignments *repository.AssignmentRepository, users *repository.UserRepository) *AssignmentService {
	return &AssignmentService{db: db, assignments: assignments, users: users}
}

func (s *AssignmentService) List(coachIDStr, studentIDStr, city, search string, offset, limit int) ([]repository.AssignmentRow, int64, error) {
	var coachID, studentID *uuid.UUID
	if coachIDStr != "" {
		id, err := uuid.Parse(coachIDStr)
		if err != nil {
			return nil, 0, apperrors.WithDetails(apperrors.ErrValidation, "coach id geçersiz")
		}
		coachID = &id
	}
	if studentIDStr != "" {
		id, err := uuid.Parse(studentIDStr)
		if err != nil {
			return nil, 0, apperrors.WithDetails(apperrors.ErrValidation, "student id geçersiz")
		}
		studentID = &id
	}
	var cityPtr *string
	if city != "" {
		cityPtr = &city
	}
	return s.assignments.ListActive(coachID, studentID, cityPtr, search, offset, limit)
}

func (s *AssignmentService) Assign(req dto.AssignRequest, by uuid.UUID) (*dto.AssignmentBulkResponse, error) {
	coachID, err := uuid.Parse(req.CoachID)
	if err != nil {
		return nil, apperrors.WithDetails(apperrors.ErrValidation, "coach id geçersiz")
	}
	coach, err := s.users.FindByID(coachID)
	if err != nil {
		return nil, err
	}
	if coach == nil || coach.Role != domain.RoleCoach {
		return nil, apperrors.WithDetails(apperrors.ErrValidation, "koç bulunamadı")
	}
	if !coach.IsActive {
		return nil, apperrors.WithDetails(apperrors.ErrValidation, "koç pasif")
	}

	resp := &dto.AssignmentBulkResponse{}
	for _, sid := range req.StudentIDs {
		studentID, err := uuid.Parse(sid)
		if err != nil {
			resp.Failed++
			resp.Errors = append(resp.Errors, sid+": id geçersiz")
			continue
		}
		st, err := s.users.FindByID(studentID)
		if err != nil {
			resp.Failed++
			resp.Errors = append(resp.Errors, sid+": "+err.Error())
			continue
		}
		if st == nil || st.Role != domain.RoleStudent {
			resp.Failed++
			resp.Errors = append(resp.Errors, sid+": öğrenci değil")
			continue
		}
		if _, err := s.assignments.Assign(coachID, studentID, &by); err != nil {
			resp.Failed++
			resp.Errors = append(resp.Errors, sid+": "+err.Error())
			continue
		}
		resp.Created++
	}
	return resp, nil
}

func (s *AssignmentService) Unassign(studentIDStr string) error {
	id, err := uuid.Parse(studentIDStr)
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "id geçersiz")
	}
	return s.assignments.Unassign(id)
}

func (s *AssignmentService) StudentsWithoutCoach(city, search string) ([]domain.User, error) {
	var cityPtr *string
	if city != "" {
		cityPtr = &city
	}
	return s.assignments.StudentsWithoutCoach(cityPtr, search)
}

func (s *AssignmentService) SetCoordinator(req dto.SetCoordinatorRequest) (int, error) {
	var coordinatorID *uuid.UUID
	if req.CoordinatorID != nil && *req.CoordinatorID != "" {
		id, err := uuid.Parse(*req.CoordinatorID)
		if err != nil {
			return 0, apperrors.WithDetails(apperrors.ErrValidation, "coordinator id geçersiz")
		}
		coord, err := s.users.FindByID(id)
		if err != nil {
			return 0, err
		}
		if coord == nil || coord.Role != domain.RoleCoordinator {
			return 0, apperrors.WithDetails(apperrors.ErrValidation, "koordinatör bulunamadı")
		}
		coordinatorID = &id
	}

	updated := 0
	for _, sid := range req.StudentIDs {
		studentID, err := uuid.Parse(sid)
		if err != nil {
			continue
		}
		res := s.db.Model(&domain.Student{}).Where("user_id = ?", studentID).
			Update("coordinator_id", coordinatorID)
		if res.Error == nil {
			updated += int(res.RowsAffected)
		}
	}
	return updated, nil
}

func (s *AssignmentService) SetParent(req dto.SetParentRequest) error {
	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		return apperrors.WithDetails(apperrors.ErrValidation, "student id geçersiz")
	}
	var parentID *uuid.UUID
	if req.ParentID != nil && *req.ParentID != "" {
		id, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return apperrors.WithDetails(apperrors.ErrValidation, "parent id geçersiz")
		}
		parent, err := s.users.FindByID(id)
		if err != nil {
			return err
		}
		if parent == nil || parent.Role != domain.RoleParent {
			return apperrors.WithDetails(apperrors.ErrValidation, "veli bulunamadı")
		}
		parentID = &id
	}
	return s.db.Model(&domain.Student{}).Where("user_id = ?", studentID).
		Update("parent_id", parentID).Error
}
