package pagination

import "github.com/gofiber/fiber/v2"

type Params struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func From(c *fiber.Ctx) Params {
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	size := c.QueryInt("pageSize", 20)
	if size < 1 {
		size = 20
	}
	if size > 200 {
		size = 200
	}
	return Params{Page: page, PageSize: size}
}

func (p Params) Offset() int { return (p.Page - 1) * p.PageSize }
func (p Params) Limit() int  { return p.PageSize }

type Page[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"totalCount"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
}

func NewPage[T any](items []T, total int64, p Params) Page[T] {
	return Page[T]{Items: items, TotalCount: total, Page: p.Page, PageSize: p.PageSize}
}
