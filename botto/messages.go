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
	Text     string
	Elements []interactiveElement
}

func (ctx *context) Send(text string) string {
	_, ts, _ := ctx.bot.client.PostMessage(ctx.channel, text, slack.NewPostMessageParameters())
	return ts
}

func (ctx *context) SendInteractive(
	text string,
	msg interactiveMessage,
	cb func(resp *slack.AttachmentActionCallback)) string {

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

	_, ts, _ := ctx.bot.client.PostMessage(ctx.channel, text, parm)

	resp := ctx.WaitAsync()

	if cb != nil {
		cb(resp)
	}

	return ts
}

func (ctx *context) Update(resp *slack.AttachmentActionCallback, msg interactiveMessage) {
	parm := slack.NewPostMessageParameters()

	attch := slack.Attachment{
		Fallback: msg.Text,
		Text:     msg.Text,
		Actions:  make([]slack.AttachmentAction, len(msg.Elements)),
	}

	for i, e := range msg.Elements {
		attch.Actions[i] = e.toAction()
	}

	parm.Attachments = []slack.Attachment{attch}

	if _, _, _, err := ctx.bot.client.SendMessage(
		ctx.channel,
		slack.MsgOptionUpdate(resp.MessageTs),
		slack.MsgOptionAttachments(parm.Attachments...),
		slack.MsgOptionPostMessageParameters(parm),
	); err != nil {
		log.Errorln(err)
	}
}

func (ctx *context) Delete(ts string) {
	ctx.bot.client.DeleteMessage(ctx.channel, ts)
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
