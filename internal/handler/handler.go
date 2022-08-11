package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/model"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/usecase"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	userUsecase   model.UserUsecase
	authUsecase   model.AuthUsecase
	absentUsecase model.AbsentUsecase
}

func New(userUsecase model.UserUsecase, authUsecase model.AuthUsecase, absentUsecase model.AbsentUsecase) *Handler {
	return &Handler{
		userUsecase:   userUsecase,
		authUsecase:   authUsecase,
		absentUsecase: absentUsecase,
	}
}

func (h *Handler) RegisterHandler() handlers.Response {
	paramLength := 5

	return func(b *gotgbot.Bot, c *ext.Context) error {
		ctx, err := h.authUsecase.CreateContextForUser(c.EffectiveSender.Id())
		if errors.Is(err, usecase.ErrInternal) {
			logrus.Error(err)
			return replyErrorMessageToUser(b, c, ErrInternal)
		}

		params, err := extractParamsFromText(c, paramLength)
		if err != nil {
			_, err := replyRegisterCommandHelp(b, c)
			return err
		}

		regisInput := &model.HimatroRegistrationInput{}
		res, err := jsonify(regisInput, params)
		if err != nil {
			logrus.Error(err)
			return replyErrorMessageToUser(b, c, err)
		}

		logIfError(json.Unmarshal([]byte(res), regisInput))

		err = h.userUsecase.Register(ctx, &model.RegistrationInput{
			ID:                   c.EffectiveSender.Id(),
			UserName:             c.EffectiveSender.Username(),
			Email:                regisInput.Email,
			Name:                 regisInput.Name,
			Password:             regisInput.Password,
			PasswordConfirmation: regisInput.PasswordConfirmation,
			InvitationCode:       regisInput.InvitationCode,
		})
		switch err {
		default:
			logrus.Error(err)
			return replyErrorMessageToUser(b, c, ErrInternal)
		case usecase.ErrExternalService:
			logrus.Warn(err)
			return replyErrorMessageToUser(b, c, ErrExternalService)
		case usecase.ErrExternalBadRequest:
			return replyErrorMessageToUser(b, c, ErrBadRequest)
		case usecase.ErrNotFound:
			return replyErrorMessageToUser(b, c, ErrNotFound)
		case usecase.ErrValidation:
			return replyErrorMessageToUser(b, c, ErrValidation)
		case nil:
			return replySuccessMessageToUser(b, c, "Registration Success")
		}

	}
}

func (h *Handler) LoginHandler() handlers.Response {
	paramLength := 1

	return func(b *gotgbot.Bot, c *ext.Context) error {
		ctx, err := h.authUsecase.CreateContextForUser(c.EffectiveSender.Id())
		if errors.Is(err, usecase.ErrInternal) {
			logrus.Error(err)
			return replyErrorMessageToUser(b, c, ErrInternal)
		}

		params, err := extractParamsFromText(c, paramLength)
		if err != nil {
			_, err := replyLoginCommandHelp(b, c)
			return err
		}

		userID := c.EffectiveSender.Id()
		password := params[0]

		_, err = h.userUsecase.LoginByPassword(ctx, userID, password)
		switch err {
		default:
			logrus.Error(err)
			return replyErrorMessageToUser(b, c, ErrInternal)
		case usecase.ErrNotFound:
			return replyErrorMessageToUser(b, c, ErrNotFound)
		case usecase.ErrValidation:
			return replyErrorMessageToUser(b, c, ErrValidation)
		case usecase.ErrExternalService:
			return replyErrorMessageToUser(b, c, ErrExternalService)
		case usecase.ErrUnauthorized:
			return replyErrorMessageToUser(b, c, ErrUnauthorized)
		case nil:
			return replySuccessMessageToUser(b, c, "login success")
		}
	}
}

func (h *Handler) GetAbsentFormHandler() handlers.Response {
	paramLength := 1

	return func(b *gotgbot.Bot, c *ext.Context) error {
		ctx, err := h.authUsecase.CreateContextForUser(c.EffectiveSender.Id())
		if errors.Is(err, usecase.ErrInternal) {
			logrus.Error(err)
			return replyErrorMessageToUser(b, c, ErrInternal)
		}

		params, err := extractParamsFromText(c, paramLength)
		if err != nil {
			_, err := replyGetAbsentFormCommandHelp(b, c)
			return err
		}

		formID, err := strconv.ParseInt(params[0], 10, 64)
		if err != nil {
			return replyErrorMessageToUser(b, c, ErrBadRequest)
		}

		absentForm, err := h.absentUsecase.GetAbsentFormByID(ctx, formID)
		switch err {
		default:
			logrus.Error(err)
			return replyErrorMessageToUser(b, c, ErrInternal)
		case usecase.ErrUnauthorized:
			return replyErrorMessageToUser(b, c, ErrUnauthorized)
		case nil:
			return replyTextMessageToUser(b, c, absentForm.ToText())
		}
	}
}

func (h *Handler) FillAbsentFormHandler() handlers.Response {
	return func(b *gotgbot.Bot, c *ext.Context) error {
		ctx, err := h.authUsecase.CreateContextForUser(c.EffectiveSender.Id())
		if err != nil {
			logrus.Error(err)
			return replyErrorMessageToUser(b, c, ErrInternal)
		}

		input := removeCommandFromText(c.Message.Text)
		if input == "" || input == "help" {
			_, err := replyFillAbsentFormCommandHelp(b, c)
			return err
		}

		params := splitInputToParams(input)
		if len(params) < 2 {
			_, err := replyFillAbsentFormCommandHelp(b, c)
			return err
		}

		status, err := model.StringToStatus(params[1])
		if err != nil {
			_, err := replyFillAbsentFormCommandHelp(b, c)
			return err
		}

		if status != model.StatusPresent && len(params) != 3 {
			_, err := replyFillAbsentFormCommandHelp(b, c)
			return err
		}

		formID, err := strconv.ParseInt(params[0], 10, 64)
		if err != nil {
			_, err := replyFillAbsentFormCommandHelp(b, c)
			return err
		}

		var execuseReason string
		if len(params) == 3 {
			execuseReason = params[2]
		}

		payload := &model.FillAbsentPayload{
			Status:        string(status),
			ExecuseReason: execuseReason,
		}

		fmt.Println(payload.Status)
		fmt.Println(payload.ExecuseReason)

		err = h.absentUsecase.FillAbsentForm(ctx, formID, payload)
		switch err {
		default:
			logrus.Warn(err)
			return replyTextMessageToUser(b, c, err.Error())
		case usecase.ErrInternal:
			logrus.Error(err)
			return replyErrorMessageToUser(b, c, ErrInternal)
		case usecase.ErrValidation:
			return replyErrorMessageToUser(b, c, ErrValidation)
		case usecase.ErrUnauthorized:
			return replyErrorMessageToUser(b, c, ErrUnauthorized)
		case nil:
			return replyTextMessageToUser(b, c, "Ok.")
		}
	}
}
