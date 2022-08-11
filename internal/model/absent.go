package model

import "context"

type AbsentUsecase interface {
	GetAbsentFormByID(ctx context.Context, id int64) (*AbsentForm, error)
	FillAbsentForm(ctx context.Context, formID int64, input *FillAbsentPayload) error
}
