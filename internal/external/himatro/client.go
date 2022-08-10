package himatro

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kumparan/go-utils"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/model"
	"github.com/sirupsen/logrus"
)

var (
	registrationUrl = "rest/members/register/"
	loginUrl        = "rest/auth/login/"
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
