package usecase

import (
	"context"
	"fmt"

	"github.com/kumparan/go-utils"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/external/himatro"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/model"
	"github.com/sirupsen/logrus"
)

type absentUsecase struct {
	himatroApi model.HimatroAPI
}

func NewAbsentUsecase(himatroApi model.HimatroAPI) model.AbsentUsecase {
	return &absentUsecase{
		himatroApi: himatroApi,
	}
}

func (u *absentUsecase) GetAbsentFormByID(ctx context.Context, id int64) (*model.AbsentForm, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.DumpIncomingContext(ctx),
		"id":  id,
	})

	absentForm, err := u.himatroApi.GetAbsentFormByID(ctx, id)
	switch err {
	default:
		logger.Error(err)
		return nil, ErrInternal
	case himatro.ErrUnauthorized:
		return nil, ErrUnauthorized
	case nil:
		return absentForm, nil
	}
}

func (u *absentUsecase) FillAbsentForm(ctx context.Context, formID int64, input *model.FillAbsentPayload) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.DumpIncomingContext(ctx),
		"formID": formID,
		"input":  utils.Dump(input),
	})

	err := u.himatroApi.FillAbsentForm(ctx, formID, input)
	switch err {
	default:
		logger.Warn(err)
		return fmt.Errorf("HIMATRO service returned error: %s", err.Error())
	case himatro.ErrValidation:
		return ErrValidation
	case himatro.ErrUnauthorized:
		return ErrUnauthorized
	case himatro.ErrInternal:
		logger.Error(err)
		return ErrInternal
	case nil:
		return nil
	}
}
