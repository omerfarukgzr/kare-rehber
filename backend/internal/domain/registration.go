package domain

import (
	"time"

	"github.com/google/uuid"
)

type Registration struct {
	ID            uuid.UUID          `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	Kind          RegistrationKind   `json:"kind" gorm:"column:kind;type:registration_kind"`
	Status        RegistrationStatus `json:"status" gorm:"column:status;type:registration_status"`
	Payload       JSONB              `json:"payload" gorm:"column:payload;type:jsonb"`
	FullName      string             `json:"fullName" gorm:"column:full_name"`
	Phone         string             `json:"phone" gorm:"column:phone"`
	Email         *string            `json:"email" gorm:"column:email"`
	City          *string            `json:"city" gorm:"column:city"`
	ReviewedBy    *uuid.UUID         `json:"reviewedBy" gorm:"column:reviewed_by;type:uuid"`
	ReviewedAt    *time.Time         `json:"reviewedAt" gorm:"column:reviewed_at"`
	ReviewNote    *string            `json:"reviewNote" gorm:"column:review_note"`
	CreatedUserID *uuid.UUID         `json:"createdUserId" gorm:"column:created_user_id;type:uuid"`
	Timestamps
}

func (Registration) TableName() string { return "registrations" }
