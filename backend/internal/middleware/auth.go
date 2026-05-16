package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/koc-luk/backend/internal/auth"
	"github.com/koc-luk/backend/internal/domain"
	apperrors "github.com/koc-luk/backend/pkg/errors"
)

const (
	CtxUserID = "user_id"
	CtxRole   = "role"
)

func JWTAuth(jwtMgr *auth.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return apperrors.ErrUnauthorized
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtMgr.Parse(tokenStr)
		if err != nil {
			return apperrors.ErrUnauthorized
		}
		c.Locals(CtxUserID, claims.UserID)
		c.Locals(CtxRole, claims.Role)
		return c.Next()
	}
}

func RequireRole(roles ...domain.Role) fiber.Handler {
	allowed := make(map[domain.Role]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals(CtxRole).(domain.Role)
		if !ok {
			return apperrors.ErrUnauthorized
		}
		if _, ok := allowed[role]; !ok {
			return apperrors.ErrForbidden
		}
		return c.Next()
	}
}

func UserID(c *fiber.Ctx) (uuid.UUID, bool) {
	id, ok := c.Locals(CtxUserID).(uuid.UUID)
	return id, ok
}

func Role(c *fiber.Ctx) (domain.Role, bool) {
	r, ok := c.Locals(CtxRole).(domain.Role)
	return r, ok
}
