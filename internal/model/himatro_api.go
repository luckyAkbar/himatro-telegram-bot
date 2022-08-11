package model

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
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

type Status string

var (
	StatusPresent        Status = "PRESENT"
	StatusExecuse        Status = "EXECUSE"
	StatusPendingPresent Status = "PENDING_PRESENT"
	StatusPendingExecuse Status = "PENDING_EXECUSE"
)

func StringToStatus(s string) (Status, error) {
	str := strings.ToUpper(s)
	switch str {
	case "PRESENT":
		return StatusPresent, nil
	case "PENDING_EXECUSE":
		return StatusPendingExecuse, nil
	case "EXECUSE":
		return StatusExecuse, nil
	case "PENDING_PRESENT":
		return StatusPendingPresent, nil
	default:
		return "", errors.New("invalid status")
	}
}

type FillAbsentPayload struct {
	Status        string `json:"status" validate:"required"`
	ExecuseReason string `json:"execuse_reason" validate:"required_unless=Status PRESENT"`
}

func (i *FillAbsentPayload) Validate() error {
	return validator.Struct(i)
}

type AbsentForm struct {
	ID                      int64          `json:"id"`
	ParticipantGroupID      int64          `json:"participant_group_id"`
	StartAt                 time.Time      `json:"start_at"`
	FinishedAt              time.Time      `json:"finished_at"`
	Title                   string         `json:"title"`
	AllowUpdateByAttendee   bool           `json:"allow_update_by_attendee"`
	AllowCreateConfirmation bool           `json:"allow_create_confirmation"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
	DeletedAt               gorm.DeletedAt `json:"deleted_at"`
	CreatedBy               int64          `json:"created_by"`
	UpdatedBy               int64          `json:"updated_by"`
	DeletedBy               null.Int       `json:"deleted_by"`
}

func (af *AbsentForm) ToText() string {
	var sb strings.Builder

	now := time.Now()
	isOpen := false

	if af.StartAt.Before(now) && af.FinishedAt.After(now) {
		isOpen = true
	}

	sb.Write([]byte(fmt.Sprintf("Title: %s\n", af.Title)))
	sb.Write([]byte(fmt.Sprintf("Start at: %s\n", af.StartAt.Format(time.RFC1123))))
	sb.Write([]byte(fmt.Sprintf("Finished at: %s\n", af.FinishedAt.Format(time.RFC1123))))
	sb.Write([]byte(fmt.Sprintf("Is still open: %v\n", isOpen)))
	sb.Write([]byte(fmt.Sprintf("Allow update: %v\n", af.AllowUpdateByAttendee)))
	sb.Write([]byte(fmt.Sprintf("Allow confirmation: %v\n", af.AllowCreateConfirmation)))

	return sb.String()
}

type HimatroAPI interface {
	Register(ctx context.Context, input *HimatroRegistrationInput) error
	Login(ctx context.Context, input *LoginPayload) (*HimatroSession, error)
	GetAbsentFormByID(ctx context.Context, id int64) (*AbsentForm, error)
	FillAbsentForm(ctx context.Context, formID int64, input *FillAbsentPayload) error
}
