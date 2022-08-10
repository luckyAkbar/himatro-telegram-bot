package model

import (
	"context"
	"time"

	"github.com/luckyAkbar/himatro-telegram-bot/internal/helper"
	"github.com/sirupsen/logrus"
)

type User struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	UserName  string    `json:"username" gorm:"not null"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) Encrypt() error {
	cryptor := helper.Cryptor()

	userName, err := cryptor.Encrypt(u.UserName)
	if err != nil {
		logrus.Error(err)
		return err
	}

	email, err := cryptor.Encrypt(u.Email)
	if err != nil {
		logrus.Error(err)
		return err
	}

	u.UserName = userName
	u.Email = email

	return nil
}

func (u *User) Decrypt() error {
	cryptor := helper.Cryptor()
	userName, err := cryptor.Decrypt(u.UserName)
	if err != nil {
		logrus.Error(err)
		return err
	}

	email, err := cryptor.Decrypt(u.Email)
	if err != nil {
		logrus.Error(err)
		return err
	}

	u.UserName = userName
	u.Email = email

	return nil
}

type RegistrationInput struct {
	ID                   int64  `json:"id" validate:"required"`
	UserName             string `json:"username" validate:"required,min=8"`
	Email                string `json:"email" validate:"required,email"`
	Name                 string `json:"name" validate:"required"`
	Password             string `json:"password" validate:"required,min=8,eqfield=PasswordConfirmation"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=8,eqfield=Password"`
	InvitationCode       string `json:"invitation_code" validate:"required"`
}

func (i *RegistrationInput) Validate() error {
	return validator.Struct(i)
}

type UserUsecase interface {
	LoginByPassword(ctx context.Context, userID int64, password string) (*Session, error)
	Register(ctx context.Context, input *RegistrationInput) error
}

type UserRepository interface {
	Register(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, userID int64) (*User, error)
}
