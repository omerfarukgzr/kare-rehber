package dto

import (
	"time"

	"github.com/koc-luk/backend/internal/domain"
)

type ThreadDTO struct {
	ID              string              `json:"id"`
	Kind            domain.ThreadKind   `json:"kind"`
	Subject         string              `json:"subject"`
	OpenedByUserID  string              `json:"openedByUserId"`
	OpenedByName    string              `json:"openedByName"`
	OpenedByRole    string              `json:"openedByRole"`
	Status          domain.ThreadStatus `json:"status"`
	UnreadCount     int                 `json:"unreadCount"`
	LastMessageAt   *time.Time          `json:"lastMessageAt"`
	LastMessageBody *string             `json:"lastMessageBody"`
	CreatedAt       time.Time           `json:"createdAt"`
}

type CreateThreadRequest struct {
	Kind    domain.ThreadKind `json:"kind" validate:"required,oneof=general feedback complaint"`
	Subject string            `json:"subject" validate:"required,min=2"`
	Body    string            `json:"body" validate:"required,min=1"`
}

type SendMessageRequest struct {
	Body string `json:"body" validate:"required,min=1"`
}

type MessageDTO struct {
	ID         string    `json:"id"`
	ThreadID   string    `json:"threadId"`
	FromUserID string    `json:"fromUserId"`
	FromName   string    `json:"fromName"`
	FromRole   string    `json:"fromRole"`
	Body       string    `json:"body"`
	ReadAt     *time.Time `json:"readAt"`
	CreatedAt  time.Time `json:"createdAt"`
}
