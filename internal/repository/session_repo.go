package repository

import (
	"context"

	"github.com/kumparan/go-utils"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type sessionRepo struct {
	db *gorm.DB
}

func NewSessionRepo(db *gorm.DB) model.SessionRepo {
	return &sessionRepo{db: db}
}

func (r *sessionRepo) Create(ctx context.Context, session *model.Session) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":     utils.DumpIncomingContext(ctx),
		"session": utils.Dump(session),
	})

	err := r.db.WithContext(ctx).Create(session).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (r *sessionRepo) FindActiveSessionByUserID(ctx context.Context, userID int64) (*model.Session, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.DumpIncomingContext(ctx),
		"userID": userID,
	})

	session := &model.Session{}
	err := r.db.WithContext(ctx).Model(&model.Session{}).
		Limit(1).Order("expired_at DESC").
		Where("user_id = ?", userID).Take(session).Error

	switch err {
	default:
		logger.Error(err)
		return nil, err
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	case nil:
		return session, nil
	}
}
