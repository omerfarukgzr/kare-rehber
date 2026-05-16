package sms

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/koc-luk/backend/internal/domain"
)

type SendInput struct {
	ToPhone     string
	ToUserID    *uuid.UUID
	Body        string
	SentBy      *uuid.UUID
	TemplateKey *string
}

type SendResult struct {
	LogID  uuid.UUID
	Status domain.SmsStatus
}

type Provider interface {
	Name() string
	Send(ctx context.Context, in SendInput) (*SendResult, error)
}

type MockProvider struct {
	db *gorm.DB
}

func NewMockProvider(db *gorm.DB) *MockProvider {
	return &MockProvider{db: db}
}

func (p *MockProvider) Name() string { return "mock" }

func (p *MockProvider) Send(ctx context.Context, in SendInput) (*SendResult, error) {
	logID := uuid.New()
	logRow := domain.SmsLog{
		ID:           logID,
		ToPhone:      in.ToPhone,
		ToUserID:     in.ToUserID,
		Body:         in.Body,
		SentByUserID: in.SentBy,
		Status:       domain.SmsMockSent,
		ProviderName: p.Name(),
		ProviderResponse: domain.JSONB{
			"mock":      true,
			"timestamp": time.Now().UTC(),
		},
		TemplateKey: in.TemplateKey,
		SentAt:      time.Now(),
	}
	if err := p.db.WithContext(ctx).Create(&logRow).Error; err != nil {
		return nil, err
	}
	return &SendResult{LogID: logID, Status: domain.SmsMockSent}, nil
}
