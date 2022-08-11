package auth

import (
	"context"

	"github.com/luckyAkbar/himatro-telegram-bot/internal/model"
)

type CtxKey string

var (
	Key CtxKey = "github.com/luckyAkbar/himatro-telegram-bot/auth/ctx"
)

func GetSessionFromContext(ctx context.Context) *model.Session {
	sess, ok := ctx.Value(Key).(*model.Session)
	if !ok {
		return nil
	}

	return sess
}
