package slackbot

import (
	"strings"

	"github.com/nlopes/slack"
)

type Context struct {
	Username string
	Channel  string
	Data     interface{}

	bot   *Bot
	sync  chan *slack.MessageEvent
	async chan *slack.AttachmentActionCallback
}

type ContextFn func(*Context, *slack.MessageEvent) error

func newContext(bot *Bot, username string, channel string) *Context {
	return &Context{
		Username: username,
		Channel:  channel,

		bot:   bot,
		sync:  make(chan *slack.MessageEvent),
		async: make(chan *slack.AttachmentActionCallback),
	}
}

func (ctx *Context) start(initFn ContextFn) error {
	if initFn != nil {
		if err := initFn(ctx, ctx.Receive()); err != nil {
			return err
		}
	}

	for {
		msg := ctx.Receive()
		ctx.exec(msg)
	}
}

func (ctx *Context) showHelp() {
	text := "Here's some stuff you can do:\n"

	ctx.bot.commands.Range(func(k, v interface{}) bool {
		cmd := v.(command)
		text += "`" + cmd.name + "` to " + cmd.description + "\n"
		return true
	})

	ctx.Send(text)
}

func (ctx *Context) exec(ev *slack.MessageEvent) {
	name := strings.TrimSpace(ev.Text)

	cmd, loaded := ctx.bot.commands.Load(name)
	if !loaded {
		ctx.showHelp()
	} else {
		cmd.(command).callback(ctx, ev)
	}
}

func (ctx *Context) dispatchSync(ev *slack.MessageEvent) {
	ctx.sync <- ev
}

func (ctx *Context) dispatchAsync(ev *slack.AttachmentActionCallback) {
	ctx.async <- ev
}

func (ctx *Context) Receive() *slack.MessageEvent {
	return <-ctx.sync
}

func (ctx *Context) ReceiveAsync() *slack.AttachmentActionCallback {
	return <-ctx.async
}
