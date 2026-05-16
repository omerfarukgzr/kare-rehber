package dto

import (
	"time"

	"github.com/koc-luk/backend/internal/domain"
)

type WeekDTO struct {
	ID            string     `json:"id"`
	WeekNo        int        `json:"weekNo"`
	Label         string     `json:"label"`
	StartDate     string     `json:"startDate"`
	EndDate       string     `json:"endDate"`
	IsOpen        bool       `json:"isOpen"`
	OpenedAt      *time.Time `json:"openedAt"`
	ClosedAt      *time.Time `json:"closedAt"`
}

func ToWeekDTO(w *domain.EvaluationWeek) WeekDTO {
	return WeekDTO{
		ID:        w.ID.String(),
		WeekNo:    w.WeekNo,
		Label:     w.Label,
		StartDate: w.StartDate.Format("2006-01-02"),
		EndDate:   w.EndDate.Format("2006-01-02"),
		IsOpen:    w.IsOpen,
		OpenedAt:  w.OpenedAt,
		ClosedAt:  w.ClosedAt,
	}
}

type CreateWeekRequest struct {
	WeekNo    int    `json:"weekNo" validate:"required,min=1"`
	Label     string `json:"label"`
	StartDate string `json:"startDate" validate:"required"`
	EndDate   string `json:"endDate" validate:"required"`
}

type UpdateWeekRequest struct {
	Label     *string `json:"label"`
	StartDate *string `json:"startDate"`
	EndDate   *string `json:"endDate"`
}

type GenerateWeeksRequest struct {
	StartDate string `json:"startDate" validate:"required"`
	WeekCount int    `json:"weekCount" validate:"required,min=1,max=52"`
}
