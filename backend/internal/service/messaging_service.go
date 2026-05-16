package service

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/koc-luk/backend/internal/domain"
	"github.com/koc-luk/backend/internal/dto"
	"github.com/koc-luk/backend/internal/repository"
	apperrors "github.com/koc-luk/backend/pkg/errors"
)

type MessagingService struct {
	db    *gorm.DB
	repo  *repository.MessagingRepository
	users *repository.UserRepository
}

func NewMessagingService(db *gorm.DB, r *repository.MessagingRepository, users *repository.UserRepository) *MessagingService {
	return &MessagingService{db: db, repo: r, users: users}
}

func (s *MessagingService) ListThreads(viewerID uuid.UUID, viewerRole domain.Role, kind, status, search string, offset, limit int) ([]dto.ThreadDTO, int64, error) {
	rows, total, err := s.repo.ListThreads(viewerID, viewerRole, kind, status, search, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	items := make([]dto.ThreadDTO, 0, len(rows))
	for _, r := range rows {
		items = append(items, dto.ThreadDTO{
			ID:              r.ID.String(),
			Kind:            r.Kind,
			Subject:         r.Subject,
			OpenedByUserID:  r.OpenedByUserID.String(),
			OpenedByName:    r.OpenedByName,
			OpenedByRole:    r.OpenedByRole,
			Status:          r.Status,
			UnreadCount:     r.UnreadCount,
			LastMessageAt:   r.LastMessageAt,
			LastMessageBody: r.LastMessageBody,
			CreatedAt:       r.CreatedAt,
		})
	}
	return items, total, nil
}

func (s *MessagingService) Create(openerID uuid.UUID, openerRole domain.Role, req dto.CreateThreadRequest) (*dto.ThreadDTO, error) {
	if openerRole != domain.RoleStudent && openerRole != domain.RoleParent && openerRole != domain.RoleAdmin && openerRole != domain.RoleCoordinator {
		return nil, apperrors.ErrForbidden
	}
	now := time.Now()
	t := &domain.MessageThread{
		ID:             uuid.New(),
		Kind:           req.Kind,
		Subject:        req.Subject,
		OpenedByUserID: openerID,
		Status:         domain.ThreadOpen,
		LastMessageAt:  &now,
	}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(t).Error; err != nil {
			return err
		}
		m := &domain.Message{
			ID:         uuid.New(),
			ThreadID:   t.ID,
			FromUserID: openerID,
			Body:       req.Body,
			CreatedAt:  now,
		}
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	d := dto.ThreadDTO{
		ID: t.ID.String(), Kind: t.Kind, Subject: t.Subject,
		OpenedByUserID: openerID.String(), Status: t.Status,
		LastMessageAt: t.LastMessageAt, CreatedAt: t.CreatedAt,
	}
	return &d, nil
}

func (s *MessagingService) Send(threadID, fromUserID uuid.UUID, viewerRole domain.Role, body string) (*dto.MessageDTO, error) {
	t, err := s.repo.FindThread(threadID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, apperrors.ErrNotFound
	}
	if !s.canAccess(t, fromUserID, viewerRole) {
		return nil, apperrors.ErrForbidden
	}
	m, err := s.repo.AppendMessage(threadID, fromUserID, body)
	if err != nil {
		return nil, err
	}
	user, _ := s.users.FindByID(fromUserID)
	role := ""
	if user != nil {
		role = string(user.Role)
	}
	name := ""
	if user != nil {
		name = user.FullName
	}
	return &dto.MessageDTO{
		ID:         m.ID.String(),
		ThreadID:   m.ThreadID.String(),
		FromUserID: m.FromUserID.String(),
		FromName:   name,
		FromRole:   role,
		Body:       m.Body,
		ReadAt:     m.ReadAt,
		CreatedAt:  m.CreatedAt,
	}, nil
}

func (s *MessagingService) Messages(threadID, viewerID uuid.UUID, viewerRole domain.Role) ([]dto.MessageDTO, error) {
	t, err := s.repo.FindThread(threadID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, apperrors.ErrNotFound
	}
	if !s.canAccess(t, viewerID, viewerRole) {
		return nil, apperrors.ErrForbidden
	}
	rows, err := s.repo.ListMessages(threadID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.MarkRead(threadID, viewerID); err != nil {
		return nil, err
	}
	items := make([]dto.MessageDTO, 0, len(rows))
	for _, r := range rows {
		items = append(items, dto.MessageDTO{
			ID:         r.ID.String(),
			ThreadID:   r.ThreadID.String(),
			FromUserID: r.FromUserID.String(),
			FromName:   r.FromName,
			FromRole:   r.FromRole,
			Body:       r.Body,
			ReadAt:     r.ReadAt,
			CreatedAt:  r.CreatedAt,
		})
	}
	return items, nil
}

func (s *MessagingService) Close(threadID, viewerID uuid.UUID, viewerRole domain.Role) error {
	t, err := s.repo.FindThread(threadID)
	if err != nil {
		return err
	}
	if t == nil {
		return apperrors.ErrNotFound
	}
	if viewerRole != domain.RoleAdmin && viewerRole != domain.RoleCoordinator && t.OpenedByUserID != viewerID {
		return apperrors.ErrForbidden
	}
	t.Status = domain.ThreadClosed
	return s.repo.UpdateThread(t)
}

func (s *MessagingService) canAccess(t *domain.MessageThread, viewerID uuid.UUID, viewerRole domain.Role) bool {
	if viewerRole == domain.RoleAdmin || viewerRole == domain.RoleCoordinator {
		return true
	}
	return t.OpenedByUserID == viewerID
}
