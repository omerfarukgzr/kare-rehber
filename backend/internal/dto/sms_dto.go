package dto

import "time"

type SmsLogDTO struct {
	ID           string    `json:"id"`
	ToPhone      string    `json:"toPhone"`
	ToUserID     *string   `json:"toUserId"`
	Body         string    `json:"body"`
	SentByUserID *string   `json:"sentByUserId"`
	Status       string    `json:"status"`
	ProviderName string    `json:"providerName"`
	TemplateKey  *string   `json:"templateKey"`
	SentAt       time.Time `json:"sentAt"`
}

type SendSmsRequest struct {
	UserIDs     []string `json:"userIds"`
	Phones      []string `json:"phones"`
	Body        string   `json:"body" validate:"required,min=1"`
	TemplateKey string   `json:"templateKey"`
}

type SendSmsResponse struct {
	Sent   int      `json:"sent"`
	Failed int      `json:"failed"`
	Errors []string `json:"errors,omitempty"`
}

type MissingCoachRow struct {
	CoachID     string `json:"coachId"`
	CoachName   string `json:"coachName"`
	CoachPhone  string `json:"coachPhone"`
	StudentID   string `json:"studentId"`
	StudentName string `json:"studentName"`
}
