package dto

import (
	"time"

	"github.com/koc-luk/backend/internal/domain"
)

type EvaluationDTO struct {
	ID              string                  `json:"id"`
	AssignmentID    string                  `json:"assignmentId"`
	WeekID          string                  `json:"weekId"`
	WeekNo          int                     `json:"weekNo"`
	WeekLabel       string                  `json:"weekLabel"`
	WeekIsOpen      bool                    `json:"weekIsOpen"`
	CoachID         string                  `json:"coachId"`
	CoachName       string                  `json:"coachName"`
	StudentID       string                  `json:"studentId"`
	StudentName     string                  `json:"studentName"`
	StudentCity     *string                 `json:"studentCity"`
	CourseStatus    *string                 `json:"courseStatus"`
	HomeworkDone    *bool                   `json:"homeworkDone"`
	Motivation      *int16                  `json:"motivation"`
	Behavior        *int16                  `json:"behavior"`
	GeneralNote     *string                 `json:"generalNote"`
	AdminOnlyNote   *string                 `json:"adminOnlyNote,omitempty"`
	Status          domain.EvaluationStatus `json:"status"`
	SubmittedBy     string                  `json:"submittedBy"`
	SubmittedAt     time.Time               `json:"submittedAt"`
	LastEditedBy    *string                 `json:"lastEditedBy,omitempty"`
	LastEditedAt    *time.Time              `json:"lastEditedAt,omitempty"`
	ApprovedBy      *string                 `json:"approvedBy,omitempty"`
	ApprovedAt      *time.Time              `json:"approvedAt,omitempty"`
}

type CreateEvaluationRequest struct {
	StudentID     string  `json:"studentId" validate:"required,uuid"`
	WeekID        string  `json:"weekId" validate:"required,uuid"`
	CourseStatus  *string `json:"courseStatus"`
	HomeworkDone  *bool   `json:"homeworkDone"`
	Motivation    *int16  `json:"motivation" validate:"omitempty,min=1,max=5"`
	Behavior      *int16  `json:"behavior" validate:"omitempty,min=1,max=5"`
	GeneralNote   *string `json:"generalNote"`
	AdminOnlyNote *string `json:"adminOnlyNote"`
}

type UpdateEvaluationRequest struct {
	CourseStatus  *string `json:"courseStatus"`
	HomeworkDone  *bool   `json:"homeworkDone"`
	Motivation    *int16  `json:"motivation" validate:"omitempty,min=1,max=5"`
	Behavior      *int16  `json:"behavior" validate:"omitempty,min=1,max=5"`
	GeneralNote   *string `json:"generalNote"`
	AdminOnlyNote *string `json:"adminOnlyNote"`
	ChangeReason  *string `json:"changeReason"`
}

type ApproveEvaluationRequest struct {
	ChangeReason *string `json:"changeReason"`
}

type EvaluationVersionDTO struct {
	ID            string       `json:"id"`
	VersionNo     int          `json:"versionNo"`
	Snapshot      domain.JSONB `json:"snapshot"`
	EditedBy      *string      `json:"editedBy"`
	EditedAt      time.Time    `json:"editedAt"`
	ChangeReason  *string      `json:"changeReason"`
}

func ToEvaluationVersionDTO(v *domain.EvaluationVersion) EvaluationVersionDTO {
	d := EvaluationVersionDTO{
		ID:           v.ID.String(),
		VersionNo:    v.VersionNo,
		Snapshot:     v.Snapshot,
		EditedAt:     v.EditedAt,
		ChangeReason: v.ChangeReason,
	}
	if v.EditedBy != nil {
		s := v.EditedBy.String()
		d.EditedBy = &s
	}
	return d
}
