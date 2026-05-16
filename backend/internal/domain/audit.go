package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID          uuid.UUID  `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	ActorUserID *uuid.UUID `json:"actorUserId" gorm:"column:actor_user_id;type:uuid"`
	Entity      string     `json:"entity" gorm:"column:entity"`
	EntityID    string     `json:"entityId" gorm:"column:entity_id"`
	Action      string     `json:"action" gorm:"column:action"`
	Before      JSONB      `json:"before" gorm:"column:before;type:jsonb"`
	After       JSONB      `json:"after" gorm:"column:after;type:jsonb"`
	Metadata    JSONB      `json:"metadata" gorm:"column:metadata;type:jsonb"`
	At          time.Time  `json:"at" gorm:"column:at"`
}

func (AuditLog) TableName() string { return "audit_logs" }
