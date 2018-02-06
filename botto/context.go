package botto

import (
	"github.com/8-team/bacotto/db"
	"github.com/nlopes/slack"
)

type userContext struct {
	user          *db.User
	currentDevice *db.Otto

	dispatcher func(bot *slackbot, ev contextEvent)
}

type contextEvent interface {
	user() string
	channel() string
	text() string
}

func (uc *userContext) init(ev contextEvent) {
	uc.user = new(db.User)

	if err := db.DB.First(uc.user, "username = ?", ev.user()).Error; err != nil {
		log.Debugln("User not found in DB, proceeding with registration")
		uc.dispatcher = uc.registerUser
	} else {
		uc.dispatcher = uc.parseCommand
	}
}

type interactiveResponse struct {
	*slack.AttachmentActionCallback
}

func (ir *interactiveResponse) user() string {
	return ir.User.Name
}

func (ir *interactiveResponse) channel() string {
	return ir.Channel.ID
}

func (ir *interactiveResponse) text() string {
	return ir.OriginalMessage.Text
}

type messageEvent struct {
	*slack.MessageEvent
}

func (ev *messageEvent) user() string {
	return ev.Msg.User
}

func (ev *messageEvent) channel() string {
	return ev.Msg.Channel
}

func (ev *messageEvent) text() string {
	return ev.Msg.Text
}
