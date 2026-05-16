package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/koc-luk/backend/internal/domain"
)

type MessagingRepository struct {
	db *gorm.DB
}

func NewMessagingRepository(db *gorm.DB) *MessagingRepository {
	return &MessagingRepository{db: db}
}

type ThreadRow struct {
	domain.MessageThread
	OpenedByName     string  `json:"openedByName" gorm:"column:opened_by_name"`
	OpenedByRole     string  `json:"openedByRole" gorm:"column:opened_by_role"`
	UnreadCount      int     `json:"unreadCount" gorm:"column:unread_count"`
	LastMessageBody  *string `json:"lastMessageBody" gorm:"column:last_message_body"`
}

func (r *MessagingRepository) ListThreads(viewerID uuid.UUID, viewerRole domain.Role, kind, status, search string, offset, limit int) ([]ThreadRow, int64, error) {
	q := r.db.Table("message_threads AS t").
		Select(`t.*,
			u.full_name AS opened_by_name,
			u.role::text AS opened_by_role,
			(SELECT COUNT(*) FROM messages m WHERE m.thread_id = t.id AND m.from_user_id != ? AND m.read_at IS NULL) AS unread_count,
			(SELECT body FROM messages m2 WHERE m2.thread_id = t.id ORDER BY m2.created_at DESC LIMIT 1) AS last_message_body`,
			viewerID,
		).
		Joins("JOIN users u ON u.id = t.opened_by_user_id")

	switch viewerRole {
	case domain.RoleAdmin, domain.RoleCoordinator:
		// admin/koordinatör tüm thread'leri görebilir (koordinatör kendi öğrencilerinin / vellerin oluşturduklarını da)
	default:
		// öğrenci/veli sadece kendi açtığı thread'leri görür
		q = q.Where("t.opened_by_user_id = ?", viewerID)
	}

	if kind != "" {
		q = q.Where("t.kind = ?", kind)
	}
	if status != "" {
		q = q.Where("t.status = ?", status)
	}
	if search != "" {
		s := "%" + search + "%"
		q = q.Where("t.subject ILIKE ? OR u.full_name ILIKE ?", s, s)
	}

	var total int64
	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []ThreadRow
	if err := q.Order("t.last_message_at DESC NULLS LAST, t.created_at DESC").
		Offset(offset).Limit(limit).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *MessagingRepository) FindThread(id uuid.UUID) (*domain.MessageThread, error) {
	var t domain.MessageThread
	err := r.db.First(&t, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *MessagingRepository) CreateThread(t *domain.MessageThread) error {
	return r.db.Create(t).Error
}

func (r *MessagingRepository) UpdateThread(t *domain.MessageThread) error {
	return r.db.Save(t).Error
}

func (r *MessagingRepository) AppendMessage(threadID uuid.UUID, fromUser uuid.UUID, body string) (*domain.Message, error) {
	now := time.Now()
	m := &domain.Message{
		ID:         uuid.New(),
		ThreadID:   threadID,
		FromUserID: fromUser,
		Body:       body,
		CreatedAt:  now,
	}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		if err := tx.Model(&domain.MessageThread{}).Where("id = ?", threadID).
			Updates(map[string]any{"last_message_at": now, "status": domain.ThreadOpen}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}

type MessageRow struct {
	domain.Message
	FromName string `json:"fromName" gorm:"column:from_name"`
	FromRole string `json:"fromRole" gorm:"column:from_role"`
}

func (r *MessagingRepository) ListMessages(threadID uuid.UUID) ([]MessageRow, error) {
	var rows []MessageRow
	err := r.db.Table("messages AS m").
		Select(`m.*, u.full_name AS from_name, u.role::text AS from_role`).
		Joins("JOIN users u ON u.id = m.from_user_id").
		Where("m.thread_id = ?", threadID).
		Order("m.created_at ASC").
		Scan(&rows).Error
	return rows, err
}

func (r *MessagingRepository) MarkRead(threadID, viewerID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&domain.Message{}).
		Where("thread_id = ? AND from_user_id != ? AND read_at IS NULL", threadID, viewerID).
		Update("read_at", now).Error
}
