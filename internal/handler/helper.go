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

// func textify(source any) string {
// 	var sb strings.Builder
// 	e := reflect.ValueOf(source).Elem()

// 	for i := 0; i < e.NumField(); i++ {
// 		field := e.Type().Field(i).Name
// 		value := e.Field(i).Interface()

// 		sb.Write([]byte(fmt.Sprintf("%s: %v\n", field, value)))
// 	}

// 	return sb.String()
// }

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

func replyGetAbsentFormCommandHelp(b *gotgbot.Bot, ctx *ext.Context) (*gotgbot.Message, error) {
	return ctx.EffectiveMessage.Reply(
		b,
		`usage: /get-absent-form <form ID>

		example: /get-absent-form 123456789
		`,
		&gotgbot.SendMessageOpts{
			ReplyToMessageId: ctx.Message.MessageId,
		},
	)
}

func replyFillAbsentFormCommandHelp(b *gotgbot.Bot, ctx *ext.Context) (*gotgbot.Message, error) {
	return ctx.EffectiveMessage.Reply(
		b,
		`usage: /fill-absent-form <form ID>, <status>, <reason>

		example: /get-absent-form 123456789, EXECUSE, I have to help my mom

		note: remember to use the commas, and don't flip the order of the input.
		also, the parameters status valid value is only one of the following:
		1. PRESENT
		2. EXECUSE
		3. PENDING PRESENT
		4. PENDING EXECUSE
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

func replySuccessMessageToUser(b *gotgbot.Bot, ctx *ext.Context, input any) error {
	// res, err := json.Marshal(input)
	// if err != nil {
	// 	logrus.Error(err)
	// 	return err
	// }

	smg := fmt.Sprintf("%+v", input)

	_, err := ctx.EffectiveMessage.Reply(
		b,
		//string(res),
		smg,
		&gotgbot.SendMessageOpts{
			ReplyToMessageId: ctx.Message.MessageId,
		},
	)

	logIfError(err)
	return err
}

func replyTextMessageToUser(b *gotgbot.Bot, ctx *ext.Context, text string) error {
	_, err := ctx.EffectiveMessage.Reply(
		b,
		text,
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
