package dto

import "time"

type AssignmentDTO struct {
	ID           string     `json:"id"`
	CoachID      string     `json:"coachId"`
	CoachName    string     `json:"coachName"`
	CoachPhone   string     `json:"coachPhone"`
	StudentID    string     `json:"studentId"`
	StudentName  string     `json:"studentName"`
	StudentPhone string     `json:"studentPhone"`
	StudentCity  *string    `json:"studentCity"`
	StartedAt    time.Time  `json:"startedAt"`
	EndedAt      *time.Time `json:"endedAt"`
	IsActive     bool       `json:"isActive"`
}

type AssignRequest struct {
	CoachID    string   `json:"coachId" validate:"required,uuid"`
	StudentIDs []string `json:"studentIds" validate:"required,min=1,dive,uuid"`
}

type UnassignRequest struct {
	StudentID string `json:"studentId" validate:"required,uuid"`
}

type AssignmentBulkResponse struct {
	Created  int      `json:"created"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors,omitempty"`
}

type SetCoordinatorRequest struct {
	StudentIDs    []string `json:"studentIds" validate:"required,min=1,dive,uuid"`
	CoordinatorID *string  `json:"coordinatorId"`
}

type SetParentRequest struct {
	StudentID string `json:"studentId" validate:"required,uuid"`
	ParentID  *string `json:"parentId"`
}
