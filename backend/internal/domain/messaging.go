package domain

import (
	"time"

	"github.com/google/uuid"
)

type MessageThread struct {
	ID                  uuid.UUID    `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	Kind                ThreadKind   `json:"kind" gorm:"column:kind;type:thread_kind"`
	Subject             string       `json:"subject" gorm:"column:subject"`
	OpenedByUserID      uuid.UUID    `json:"openedByUserId" gorm:"column:opened_by_user_id;type:uuid"`
	AssignedToUserID    *uuid.UUID   `json:"assignedToUserId" gorm:"column:assigned_to_user_id;type:uuid"`
	Status              ThreadStatus `json:"status" gorm:"column:status;type:thread_status"`
	LastMessageAt       *time.Time   `json:"lastMessageAt" gorm:"column:last_message_at"`
	Timestamps

	OpenedBy *User      `json:"openedBy,omitempty" gorm:"foreignKey:OpenedByUserID;references:ID"`
	Messages []Message  `json:"messages,omitempty" gorm:"foreignKey:ThreadID;references:ID"`
}

func (MessageThread) TableName() string { return "message_threads" }

type Message struct {
	ID         uuid.UUID  `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	ThreadID   uuid.UUID  `json:"threadId" gorm:"column:thread_id;type:uuid"`
	FromUserID uuid.UUID  `json:"fromUserId" gorm:"column:from_user_id;type:uuid"`
	Body       string     `json:"body" gorm:"column:body"`
	ReadAt     *time.Time `json:"readAt" gorm:"column:read_at"`
	CreatedAt  time.Time  `json:"createdAt" gorm:"column:created_at"`

	From *User `json:"from,omitempty" gorm:"foreignKey:FromUserID;references:ID"`
}

func (Message) TableName() string { return "messages" }
