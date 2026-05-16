package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/koc-luk/backend/internal/domain"
)

type Claims struct {
	UserID uuid.UUID   `json:"sub"`
	Role   domain.Role `json:"role"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret    []byte
	expiresIn time.Duration
}

func NewManager(secret string, expiresIn time.Duration) *Manager {
	return &Manager{secret: []byte(secret), expiresIn: expiresIn}
}

func (m *Manager) Generate(userID uuid.UUID, role domain.Role) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(m.expiresIn)
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			Subject:   userID.String(),
			Issuer:    "kare-rehber",
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, exp, nil
}

func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
