package domain

import (
	"time"

	"github.com/google/uuid"
)

type SmsLog struct {
	ID               uuid.UUID  `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	ToPhone          string     `json:"toPhone" gorm:"column:to_phone"`
	ToUserID         *uuid.UUID `json:"toUserId" gorm:"column:to_user_id;type:uuid"`
	Body             string     `json:"body" gorm:"column:body"`
	SentByUserID     *uuid.UUID `json:"sentByUserId" gorm:"column:sent_by_user_id;type:uuid"`
	Status           SmsStatus  `json:"status" gorm:"column:status;type:sms_status"`
	ProviderName     string     `json:"providerName" gorm:"column:provider_name"`
	ProviderResponse JSONB      `json:"providerResponse" gorm:"column:provider_response;type:jsonb"`
	TemplateKey      *string    `json:"templateKey" gorm:"column:template_key"`
	SentAt           time.Time  `json:"sentAt" gorm:"column:sent_at"`
}

func (SmsLog) TableName() string { return "sms_logs" }
