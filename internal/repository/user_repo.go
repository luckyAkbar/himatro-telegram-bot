package repository

import (
	"context"

	"github.com/kumparan/go-utils"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Register(ctx context.Context, user *model.User) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":  utils.DumpIncomingContext(ctx),
		"user": utils.Dump(user),
	})

	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.DumpIncomingContext(ctx),
		"id":  id,
	})

	user := &model.User{}
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(user).Error
	switch err {
	default:
		logger.Error(err)
		return nil, err
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	case nil:
		return user, nil
	}
}
