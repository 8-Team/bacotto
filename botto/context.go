package botto

import (
	"github.com/8-team/bacotto/db"
	"github.com/nlopes/slack"
)

type context struct {
	user    *db.User
	channel string
	bot     *slackbot
	sync    chan *slack.MessageEvent
	async   chan *slack.AttachmentActionCallback
}

func newContext(bot *slackbot, username string, channel string) *context {
	ctx := &context{
		user:    new(db.User),
		channel: channel,
		bot:     bot,
		sync:    make(chan *slack.MessageEvent),
		async:   make(chan *slack.AttachmentActionCallback),
	}

	ctx.user.Username = username
	return ctx
}

func (ctx *context) start() error {
	var registrationRequired bool

	if err := db.DB.Preload("Otto").First(ctx.user, "username = ?", ctx.user.Username).Error; err != nil {
		log.Debugln("User not found in DB, proceeding with registration")
		registrationRequired = true
	}

	for {
		ev := ctx.Wait()

		if registrationRequired {
			ctx.registerUser(ev)
			registrationRequired = false
		} else {
			ctx.parseCommand(ev)
		}
	}
}

func (ctx *context) dispatchSync(ev *slack.MessageEvent) {
	ctx.sync <- ev
}

func (ctx *context) dispatchAsync(ev *slack.AttachmentActionCallback) {
	ctx.async <- ev
}

func (ctx *context) Wait() *slack.MessageEvent {
	return <-ctx.sync
}

func (ctx *context) WaitAsync() *slack.AttachmentActionCallback {
	return <-ctx.async
}
