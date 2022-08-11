package himatro

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/kumparan/go-utils"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/auth"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/config"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/model"
	"github.com/sirupsen/logrus"
)

var (
	registrationUrl   = "rest/members/register/"
	loginUrl          = "rest/auth/login/"
	getAbsentFormUrl  = "rest/absent/form/"
	fillAbsentFormUrl = "rest/absent/form/"
)

type himatroApi struct {
	host string
}

func NewClient(host string) model.HimatroAPI {
	return &himatroApi{
		host: host,
	}
}

func (h *himatroApi) Register(ctx context.Context, input *model.HimatroRegistrationInput) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"input": utils.Dump(input),
	})

	if err := input.Validate(); err != nil {
		return ErrValidation
	}

	payload, err := json.Marshal(input)
	if err != nil {
		logger.Error(err)
		return ErrInternal
	}

	url := fmt.Sprintf("%s%s", h.host, registrationUrl)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		logger.Error(err)
		return ErrInternal
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	default:
		logger.Error(err)
		return ErrExternalService
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusOK:
		return nil
	}
}

func (h *himatroApi) Login(ctx context.Context, input *model.LoginPayload) (*model.HimatroSession, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"email": input.Email,
	})

	if err := input.Validate(); err != nil {
		return nil, ErrValidation
	}

	payload, err := json.Marshal(input)
	if err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	url := fmt.Sprintf("%s%s", h.host, loginUrl)
	resp, err := http.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	default:
		logger.Error(err)
		return nil, ErrExternalService
	case http.StatusNotFound:
		return nil, ErrNotFound
	case http.StatusUnauthorized:
		return nil, ErrUnauthorized
	case http.StatusOK:
		break
	}

	session := &model.HimatroSession{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	err = json.Unmarshal(body, session)
	if err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	return session, nil
}

func (h *himatroApi) GetAbsentFormByID(ctx context.Context, id int64) (*model.AbsentForm, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.DumpIncomingContext(ctx),
		"id":  id,
	})

	sess := auth.GetSessionFromContext(ctx)
	if sess == nil {
		return nil, ErrUnauthorized
	}

	url := fmt.Sprintf("%s%s%d/", config.HimatroAPIHost(), getAbsentFormUrl, id)
	req, err := h.createAuthorizedReq(http.MethodGet, url, nil, sess, "")
	if err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	defer resp.Body.Close()

	absentForm := &model.AbsentForm{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	err = json.Unmarshal(body, absentForm)
	if err != nil {
		logger.Error(err)
		return nil, ErrInternal
	}

	return absentForm, nil
}

func (h *himatroApi) FillAbsentForm(ctx context.Context, formID int64, input *model.FillAbsentPayload) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.DumpIncomingContext(ctx),
		"formID": formID,
		"input":  utils.Dump(input),
	})

	sess := auth.GetSessionFromContext(ctx)
	if sess == nil {
		return ErrUnauthorized
	}

	if err := input.Validate(); err != nil {
		return ErrValidation
	}

	payload, err := json.Marshal(input)
	if err != nil {
		logger.Error(err)
		return ErrInternal
	}

	url := fmt.Sprintf("%s%s%d/", h.host, fillAbsentFormUrl, formID)
	req, err := h.createAuthorizedReq(http.MethodPost, url, bytes.NewBuffer(payload), sess, "application/json")
	if err != nil {
		logger.Error(err)
		return ErrInternal
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return ErrInternal
	}

	defer resp.Body.Close()

	type message struct {
		Message string `json:"message"`
	}

	msg := &message{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return ErrInternal
	}

	err = json.Unmarshal(body, msg)
	if err != nil {
		logger.Error(err)
		return ErrInternal
	}

	if resp.StatusCode != http.StatusOK {
		logger.Warn(msg.Message)

		return errors.New(msg.Message)
	}

	return nil
}

func (h *himatroApi) createAuthorizedReq(method, url string, body io.Reader, session *model.Session, contentType string) (*http.Request, error) {
	bearer := fmt.Sprintf("Bearer %s", session.AccessToken)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	req.Header.Add("Authorization", bearer)

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	return req, nil
}
