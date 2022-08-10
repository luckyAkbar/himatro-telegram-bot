package usecase

import (
	"context"
	"time"

	"github.com/kumparan/go-utils"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/external/himatro"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/model"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/repository"
	"github.com/sirupsen/logrus"
)

type userUsecase struct {
	userRepo      model.UserRepository
	himatroClient model.HimatroAPI
	sessionRepo   model.SessionRepo
}

func NewUserUsecase(userRepo model.UserRepository, himatroClient model.HimatroAPI, sessionRepo model.SessionRepo) model.UserUsecase {
	return &userUsecase{
		userRepo:      userRepo,
		himatroClient: himatroClient,
		sessionRepo:   sessionRepo,
	}
}

func (u *userUsecase) LoginByPassword(ctx context.Context, userID int64, password string) (*model.Session, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.DumpIncomingContext(ctx),
		"userID": userID,
	})

	user, err := u.userRepo.GetUserByID(ctx, userID)
	switch err {
	default:
		logger.Error(err)
		return nil, ErrInternal
	case repository.ErrNotFound:
		return nil, ErrNotFound
	case nil:
		break
	}

	if err := user.Decrypt(); err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	payload := &model.LoginPayload{
		Email:    user.Email,
		Password: password,
	}

	himatroSession, err := u.himatroClient.Login(ctx, payload)
	switch err {
	default:
		logger.Error(err)
		return nil, ErrInternal
	case himatro.ErrValidation:
		return nil, ErrValidation
	case himatro.ErrExternalService:
		return nil, ErrExternalService
	case himatro.ErrNotFound:
		return nil, ErrNotFound
	case himatro.ErrUnauthorized:
		return nil, ErrUnauthorized
	case nil:
		break
	}

	session := &model.Session{
		ID:           utils.GenerateID(),
		UserID:       userID,
		AccessToken:  himatroSession.AccessToken,
		RefreshToken: himatroSession.RefreshToken,
		ExpiredAt:    himatroSession.AccessTokenExpiredAt,
	}

	if err := u.sessionRepo.Create(ctx, session); err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	return session, nil
}

func (u *userUsecase) Register(ctx context.Context, input *model.RegistrationInput) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"input": utils.Dump(input),
	})

	if err := input.Validate(); err != nil {
		return ErrValidation
	}

	user := &model.User{
		ID:        input.ID,
		UserName:  input.UserName,
		Email:     input.Email,
		CreatedAt: time.Now(),
	}

	if err := user.Encrypt(); err != nil {
		logger.Error(err)
		return ErrInternal
	}

	himatroRegisInput := &model.HimatroRegistrationInput{
		Email:                input.Email,
		Name:                 input.Name,
		Password:             input.Password,
		PasswordConfirmation: input.PasswordConfirmation,
		InvitationCode:       input.InvitationCode,
	}

	err := u.himatroClient.Register(ctx, himatroRegisInput)
	switch err {
	default:
		logger.Error(err)
		return ErrInternal
	case himatro.ErrExternalService:
		logger.Warn(err)
		return ErrExternalService
	case himatro.ErrBadRequest:
		return ErrExternalBadRequest
	case himatro.ErrNotFound:
		return ErrNotFound
	case nil:
		break
	}

	if err := u.userRepo.Register(ctx, user); err != nil {
		logger.Error(err)
		return ErrInternal
	}

	return nil
}
