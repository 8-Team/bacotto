package botto

import (
	"github.com/nlopes/slack"
	uuid "github.com/satori/go.uuid"
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

type interactiveMessage struct {
	Callback eventCallback
	Text     string
	Elements []interactiveElement
}

func (b *slackbot) Message(channel string, text string) string {
	_, ts, _ := b.client.PostMessage(channel, text, slack.NewPostMessageParameters())
	return ts
}

func (b *slackbot) InteractiveMessage(channel string, text string, msg interactiveMessage) string {
	parm := slack.NewPostMessageParameters()
	uid := uuid.NewV4().String()

	attch := slack.Attachment{
		Fallback:   text,
		CallbackID: uid,
		Text:       msg.Text,
		Actions:    make([]slack.AttachmentAction, len(msg.Elements)),
	}

	for i, e := range msg.Elements {
		attch.Actions[i] = e.toAction()
	}

	parm.Attachments = []slack.Attachment{attch}

	bot.registerCallback(uid, msg.Callback)

	_, ts, _ := bot.client.PostMessage(channel, text, parm)
	return ts
}

func (b *slackbot) Update(resp *interactiveResponse, msg interactiveMessage) {
	parm := slack.NewPostMessageParameters()

	attch := slack.Attachment{
		Fallback:   resp.text(),
		CallbackID: resp.CallbackID,
		Text:       msg.Text,
		Actions:    make([]slack.AttachmentAction, len(msg.Elements)),
	}

	for i, e := range msg.Elements {
		attch.Actions[i] = e.toAction()
	}

	parm.Attachments = []slack.Attachment{attch}

	bot.client.SendMessage(
		resp.channel(),
		slack.MsgOptionUpdate(resp.MessageTs),
		slack.MsgOptionAttachments(parm.Attachments...),
		slack.MsgOptionPostMessageParameters(parm),
	)
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
