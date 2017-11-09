package botto

import (
	"github.com/nlopes/slack"
)

type interactiveElement interface {
	toAction() slack.AttachmentAction
}

type messageButton struct {
	Name  string
	Text  string
	Value string
}

type messageMenu struct {
	Name   string
	Text   string
	Values map[string]string
}

type messageFormat struct {
	Callback string
	Elements []interactiveElement
}

func (b *slackbot) Message(channel string, text string) {
	b.client.PostMessage(channel, text, slack.NewPostMessageParameters())
}

func (b *slackbot) InteractiveMessage(channel string, text string, msg messageFormat) {
	parm := slack.NewPostMessageParameters()

	attch := slack.Attachment{
		Fallback:   text,
		CallbackID: msg.Callback,
		Actions:    make([]slack.AttachmentAction, len(msg.Elements)),
	}

	for i, e := range msg.Elements {
		attch.Actions[i] = e.toAction()
	}

	parm.Attachments = []slack.Attachment{attch}

	bot.client.PostMessage(channel, text, parm)
}

func (mb messageButton) toAction() slack.AttachmentAction {
	return slack.AttachmentAction{
		Name:  mb.Name,
		Text:  mb.Text,
		Type:  "button",
		Value: mb.Value,
	}
}

func (mm messageMenu) toAction() slack.AttachmentAction {
	opts := make([]slack.AttachmentActionOption, 0, len(mm.Values))
	for value, name := range mm.Values {
		opts = append(opts, slack.AttachmentActionOption{
			Value: value,
			Text:  name,
		})
	}

	return slack.AttachmentAction{
		Name:    mm.Name,
		Text:    mm.Text,
		Type:    "select",
		Options: opts,
	}
}
