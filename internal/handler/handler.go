package handler

import (
	"context"
	"encoding/json"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/model"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/usecase"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	userUsecase model.UserUsecase
}

func New(userUsecase model.UserUsecase) *Handler {
	return &Handler{
		userUsecase: userUsecase,
	}
}

func (h *Handler) RegisterHandler() handlers.Response {
	paramLength := 5

	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		params, err := extractParamsFromText(ctx, paramLength)
		if err != nil {
			_, err := replyRegisterCommandHelp(b, ctx)
			return err
		}

		regisInput := &model.HimatroRegistrationInput{}
		res, err := jsonify(regisInput, params)
		if err != nil {
			logrus.Error(err)
			return replyErrorMessageToUser(b, ctx, err)
		}

		logIfError(json.Unmarshal([]byte(res), regisInput))

		err = h.userUsecase.Register(context.Background(), &model.RegistrationInput{
			ID:                   ctx.EffectiveSender.Id(),
			UserName:             ctx.EffectiveSender.Username(),
			Email:                regisInput.Email,
			Name:                 regisInput.Name,
			Password:             regisInput.Password,
			PasswordConfirmation: regisInput.PasswordConfirmation,
			InvitationCode:       regisInput.InvitationCode,
		})
		switch err {
		default:
			logrus.Error(err)
			return replyErrorMessageToUser(b, ctx, ErrInternal)
		case usecase.ErrExternalService:
			logrus.Warn(err)
			return replyErrorMessageToUser(b, ctx, ErrExternalService)
		case usecase.ErrExternalBadRequest:
			return replyErrorMessageToUser(b, ctx, ErrBadRequest)
		case usecase.ErrNotFound:
			return replyErrorMessageToUser(b, ctx, ErrNotFound)
		case usecase.ErrValidation:
			return replyErrorMessageToUser(b, ctx, ErrValidation)
		case nil:
			return replySuccessMessageToUser(b, ctx, "Registration Success")
		}

	}
}

func (h *Handler) LoginHandler() handlers.Response {
	paramLength := 1

	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		params, err := extractParamsFromText(ctx, paramLength)
		if err != nil {
			_, err := replyLoginCommandHelp(b, ctx)
			return err
		}

		userID := ctx.EffectiveSender.Id()
		password := params[0]

		_, err = h.userUsecase.LoginByPassword(context.TODO(), userID, password)
		switch err {
		default:
			logrus.Error(err)
			return replyErrorMessageToUser(b, ctx, ErrInternal)
		case usecase.ErrNotFound:
			return replyErrorMessageToUser(b, ctx, ErrNotFound)
		case usecase.ErrValidation:
			return replyErrorMessageToUser(b, ctx, ErrValidation)
		case usecase.ErrExternalService:
			return replyErrorMessageToUser(b, ctx, ErrExternalService)
		case usecase.ErrUnauthorized:
			return replyErrorMessageToUser(b, ctx, ErrUnauthorized)
		case nil:
			return replySuccessMessageToUser(b, ctx, "login success")
		}
	}
}
