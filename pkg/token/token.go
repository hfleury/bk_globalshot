package token

import (
	"time"
)

type Maker interface {
	CreateToken(email string, role string, companyID string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

type Payload struct {
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CompanyID string    `json:"company_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
