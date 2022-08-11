package model

import (
	"context"
	"time"
)

type Session struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiredAt    time.Time `json:"expired_at"`
}

type SessionRepo interface {
	Create(ctx context.Context, session *Session) error
	FindActiveSessionByUserID(ctx context.Context, userID int64) (*Session, error)
	// FindByAccessToken(ctx context.Context, accessToken string) (*Session, error)
}
