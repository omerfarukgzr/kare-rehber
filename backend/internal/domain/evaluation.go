package domain

import (
	"time"

	"github.com/google/uuid"
)

type EvaluationWeek struct {
	ID            uuid.UUID  `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	WeekNo        int        `json:"weekNo" gorm:"column:week_no"`
	Label         string     `json:"label" gorm:"column:label"`
	StartDate     time.Time  `json:"startDate" gorm:"column:start_date;type:date"`
	EndDate       time.Time  `json:"endDate" gorm:"column:end_date;type:date"`
	IsOpen        bool       `json:"isOpen" gorm:"column:is_open"`
	OpenedByUserID *uuid.UUID `json:"openedByUserId" gorm:"column:opened_by_user_id;type:uuid"`
	OpenedAt      *time.Time `json:"openedAt" gorm:"column:opened_at"`
	ClosedAt      *time.Time `json:"closedAt" gorm:"column:closed_at"`
	Timestamps
}

func (EvaluationWeek) TableName() string { return "evaluation_weeks" }

type Evaluation struct {
	ID                 uuid.UUID        `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	AssignmentID       uuid.UUID        `json:"assignmentId" gorm:"column:assignment_id;type:uuid"`
	EvaluationWeekID   uuid.UUID        `json:"evaluationWeekId" gorm:"column:evaluation_week_id;type:uuid"`
	CourseStatus       *string          `json:"courseStatus" gorm:"column:course_status"`
	HomeworkDone       *bool            `json:"homeworkDone" gorm:"column:homework_done"`
	Motivation         *int16           `json:"motivation" gorm:"column:motivation"`
	Behavior           *int16           `json:"behavior" gorm:"column:behavior"`
	GeneralNote        *string          `json:"generalNote" gorm:"column:general_note"`
	AdminOnlyNote      *string          `json:"adminOnlyNote,omitempty" gorm:"column:admin_only_note"`
	Status             EvaluationStatus `json:"status" gorm:"column:status;type:evaluation_status"`
	SubmittedBy        uuid.UUID        `json:"submittedBy" gorm:"column:submitted_by;type:uuid"`
	SubmittedAt        time.Time        `json:"submittedAt" gorm:"column:submitted_at"`
	LastEditedBy       *uuid.UUID       `json:"lastEditedBy" gorm:"column:last_edited_by;type:uuid"`
	LastEditedAt       *time.Time       `json:"lastEditedAt" gorm:"column:last_edited_at"`
	ApprovedBy         *uuid.UUID       `json:"approvedBy" gorm:"column:approved_by;type:uuid"`
	ApprovedAt         *time.Time       `json:"approvedAt" gorm:"column:approved_at"`
	Timestamps

	Assignment *CoachStudentAssignment `json:"assignment,omitempty" gorm:"foreignKey:AssignmentID;references:ID"`
	Week       *EvaluationWeek         `json:"week,omitempty" gorm:"foreignKey:EvaluationWeekID;references:ID"`
}

func (Evaluation) TableName() string { return "evaluations" }

type EvaluationVersion struct {
	ID            uuid.UUID  `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	EvaluationID  uuid.UUID  `json:"evaluationId" gorm:"column:evaluation_id;type:uuid"`
	VersionNo     int        `json:"versionNo" gorm:"column:version_no"`
	Snapshot      JSONB      `json:"snapshot" gorm:"column:snapshot;type:jsonb"`
	EditedBy      *uuid.UUID `json:"editedBy" gorm:"column:edited_by;type:uuid"`
	EditedAt      time.Time  `json:"editedAt" gorm:"column:edited_at"`
	ChangeReason  *string    `json:"changeReason" gorm:"column:change_reason"`
}

func (EvaluationVersion) TableName() string { return "evaluation_versions" }
