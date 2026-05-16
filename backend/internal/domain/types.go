package domain

import "time"

type Role string

const (
	RoleAdmin       Role = "admin"
	RoleCoach       Role = "coach"
	RoleStudent     Role = "student"
	RoleParent      Role = "parent"
	RoleCoordinator Role = "coordinator"
)

func (r Role) Valid() bool {
	switch r {
	case RoleAdmin, RoleCoach, RoleStudent, RoleParent, RoleCoordinator:
		return true
	}
	return false
}

type RegistrationStatus string

const (
	RegistrationPending  RegistrationStatus = "pending"
	RegistrationApproved RegistrationStatus = "approved"
	RegistrationRejected RegistrationStatus = "rejected"
)

type RegistrationKind string

const (
	RegistrationKindStudent RegistrationKind = "student"
	RegistrationKindCoach   RegistrationKind = "coach"
)

type CoachStatus string

const (
	CoachPending   CoachStatus = "pending"
	CoachApproved  CoachStatus = "approved"
	CoachRejected  CoachStatus = "rejected"
	CoachSuspended CoachStatus = "suspended"
)

type EvaluationStatus string

const (
	EvalPending       EvaluationStatus = "pending"
	EvalApproved      EvaluationStatus = "approved"
	EvalEditedByAdmin EvaluationStatus = "edited_by_admin"
)

type ThreadKind string

const (
	ThreadGeneral   ThreadKind = "general"
	ThreadFeedback  ThreadKind = "feedback"
	ThreadComplaint ThreadKind = "complaint"
)

type ThreadStatus string

const (
	ThreadOpen   ThreadStatus = "open"
	ThreadClosed ThreadStatus = "closed"
)

type SmsStatus string

const (
	SmsMockSent SmsStatus = "mock_sent"
	SmsSent     SmsStatus = "sent"
	SmsFailed   SmsStatus = "failed"
)

// BaseModel is embedded by tables with id+timestamps
type Timestamps struct {
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}
