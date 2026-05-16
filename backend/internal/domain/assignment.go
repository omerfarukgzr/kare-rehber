package domain

import (
	"time"

	"github.com/google/uuid"
)

type CoachStudentAssignment struct {
	ID         uuid.UUID  `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	CoachID    uuid.UUID  `json:"coachId" gorm:"column:coach_id;type:uuid"`
	StudentID  uuid.UUID  `json:"studentId" gorm:"column:student_id;type:uuid"`
	StartedAt  time.Time  `json:"startedAt" gorm:"column:started_at"`
	EndedAt    *time.Time `json:"endedAt" gorm:"column:ended_at"`
	IsActive   bool       `json:"isActive" gorm:"column:is_active"`
	AssignedBy *uuid.UUID `json:"assignedBy" gorm:"column:assigned_by;type:uuid"`
	Timestamps

	Coach   *User `json:"coach,omitempty" gorm:"foreignKey:CoachID;references:ID"`
	Student *User `json:"student,omitempty" gorm:"foreignKey:StudentID;references:ID"`
}

func (CoachStudentAssignment) TableName() string { return "coach_student_assignments" }
