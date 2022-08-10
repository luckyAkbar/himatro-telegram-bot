package model

import (
	"context"
	"time"
)

type HimatroRegistrationInput struct {
	Email                string `json:"email" validate:"required,email"`
	Name                 string `json:"name" validate:"required"`
	Password             string `json:"password" validate:"required,min=8,eqfield=PasswordConfirmation"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=8,eqfield=Password"`
	InvitationCode       string `json:"invitation_code" validate:"required"`
}

func (i *HimatroRegistrationInput) Validate() error {
	return validator.Struct(i)
}

type HimatroSession struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiredAt  time.Time `json:"access_token_expired_at"`
	RefreshTokenExpiredAt time.Time `json:"refresh_token_expired_at"`
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (i *LoginPayload) Validate() error {
	return validator.Struct(i)
}

type HimatroAPI interface {
	Register(ctx context.Context, input *HimatroRegistrationInput) error
	Login(ctx context.Context, input *LoginPayload) (*HimatroSession, error)
}
