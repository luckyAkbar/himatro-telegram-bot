package model

import "context"

type AuthUsecase interface {
	CreateContextForUser(userID int64) (context.Context, error)
}
