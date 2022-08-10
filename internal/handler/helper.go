package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/sirupsen/logrus"
)

func removeCommandFromText(s string) string {
	parts := strings.Split(s, " ")
	if len(parts) < 2 {
		return ""
	}

	pure := strings.Join(parts[1:], " ")

	return pure
}

func splitInputToParams(input string) []string {
	params := strings.Split(input, ",")
	for i, param := range params {
		params[i] = strings.Trim(param, " ")
	}

	return params
}

// func sendGeneralClueWhenInvalidCommand(b *gotgbot.Bot, ctx *ext.Context) (*gotgbot.Message, error) {
// 	return ctx.EffectiveMessage.Reply(b, "invalid params. Please type help after your command to get help.", &gotgbot.SendMessageOpts{
// 		ReplyToMessageId: ctx.Message.MessageId,
// 	})
// }

func jsonify[T any](target T, data []string) (string, error) {
	var temp = map[string]string{}
	e := reflect.ValueOf(target).Elem()

	if e.NumField() != len(data) {
		return "", errors.New("parsing error: not enough data")
	}

	for i := 0; i < e.NumField(); i++ {
		field, ok := reflect.TypeOf(target).Elem().FieldByName(e.Type().Field(i).Name)
		if !ok {
			continue
		}

		name := string(field.Tag.Get("json"))
		temp[name] = data[i]
	}

	res, err := json.Marshal(temp)
	return string(res), err
}

func replyRegisterCommandHelp(b *gotgbot.Bot, ctx *ext.Context) (*gotgbot.Message, error) {
	return ctx.EffectiveMessage.Reply(
		b,
		`usage: /register <your@email.address>, <your full name>, <your password>, <your confirmation password>, <your invitation code>

		example: /register lucky@akbar.tech, lucky akbar, password, password, 123456789
		note: remember to use the commas, and don't flip the order of the input.
		`,
		&gotgbot.SendMessageOpts{
			ReplyToMessageId: ctx.Message.MessageId,
		},
	)
}

func replyLoginCommandHelp(b *gotgbot.Bot, ctx *ext.Context) (*gotgbot.Message, error) {
	return ctx.EffectiveMessage.Reply(
		b,
		`usage: /login <your password>

		example: /login mysuperSecretDontuseThis
		`,
		&gotgbot.SendMessageOpts{
			ReplyToMessageId: ctx.Message.MessageId,
		},
	)
}

func replyErrorMessageToUser(b *gotgbot.Bot, ctx *ext.Context, actualErr error) error {
	_, err := ctx.EffectiveMessage.Reply(
		b,
		fmt.Sprintf("An error occurred: %s", actualErr.Error()),
		&gotgbot.SendMessageOpts{
			ReplyToMessageId: ctx.Message.MessageId,
		},
	)

	logIfError(err)

	return actualErr
}

func replySuccessMessageToUser(b *gotgbot.Bot, ctx *ext.Context, msg string) error {
	_, err := ctx.EffectiveMessage.Reply(
		b,
		msg,
		&gotgbot.SendMessageOpts{
			ReplyToMessageId: ctx.Message.MessageId,
		},
	)

	logIfError(err)
	return err
}

func logIfError(err error) {
	if err != nil {
		logrus.Error(err)
	}
}

func extractParamsFromText(ctx *ext.Context, paramLength int) ([]string, error) {
	var params []string
	text := ctx.Message.Text
	extractionErr := errors.New("extraction error")

	input := removeCommandFromText(text)
	if input == "help" {
		return params, extractionErr
	}

	if paramLength != 0 && input == "" {
		return params, extractionErr
	}

	params = splitInputToParams(input)
	if len(params) != paramLength {
		return params, extractionErr
	}

	return params, nil
}
