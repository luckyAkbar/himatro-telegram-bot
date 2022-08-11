package usecase

import (
	"context"

	"github.com/luckyAkbar/himatro-telegram-bot/internal/auth"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/model"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/repository"
	"github.com/sirupsen/logrus"
)

type authUsecase struct {
	sessionRepo model.SessionRepo
}

func NewAuthUsecase(sessionRepo model.SessionRepo) model.AuthUsecase {
	return &authUsecase{
		sessionRepo: sessionRepo,
	}
}

func (u *authUsecase) CreateContextForUser(userID int64) (context.Context, error) {
	logger := logrus.WithFields(logrus.Fields{
		"user_id": userID,
	})

	session, err := u.sessionRepo.FindActiveSessionByUserID(context.Background(), userID)
	switch err {
	default:
		logger.Error(err)
		return nil, ErrInternal
	case repository.ErrNotFound:
		return nil, ErrNotFound
	case nil:
		break
	}

	ctx := u.createContextFromSession(session)
	return ctx, nil
}

func (u *authUsecase) createContextFromSession(session *model.Session) context.Context {
	return context.WithValue(context.TODO(), auth.Key, session)
}
