package slackbot

import (
	"github.com/nlopes/slack"
	uuid "github.com/satori/go.uuid"
)

type InteractiveElement interface {
	ToAction() slack.AttachmentAction
}

type MessageButton struct {
	Name  string
	Text  string
	Value string
}

type MessageMenu struct {
	Name   string
	Text   string
	Values map[string]string
}

type InteractiveMessage struct {
	Text     string
	Elements []InteractiveElement
}

func (ctx *Context) Send(text string) string {
	_, ts, _ := ctx.bot.client.PostMessage(ctx.Channel, text, slack.NewPostMessageParameters())
	return ts
}

func (ctx *Context) SendInteractive(
	text string,
	msg InteractiveMessage,
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
		attch.Actions[i] = e.ToAction()
	}

	parm.Attachments = []slack.Attachment{attch}

	_, ts, _ := ctx.bot.client.PostMessage(ctx.Channel, text, parm)

	resp := ctx.ReceiveAsync()

	if cb != nil {
		cb(resp)
	}

	return ts
}

func (ctx *Context) Update(resp *slack.AttachmentActionCallback, msg InteractiveMessage) {
	parm := slack.NewPostMessageParameters()

	attch := slack.Attachment{
		Fallback: msg.Text,
		Text:     msg.Text,
		Actions:  make([]slack.AttachmentAction, len(msg.Elements)),
	}

	for i, e := range msg.Elements {
		attch.Actions[i] = e.ToAction()
	}

	parm.Attachments = []slack.Attachment{attch}

	if _, _, _, err := ctx.bot.client.SendMessage(
		ctx.Channel,
		slack.MsgOptionUpdate(resp.MessageTs),
		slack.MsgOptionAttachments(parm.Attachments...),
		slack.MsgOptionPostMessageParameters(parm),
	); err != nil {
		ctx.bot.logger.Errorln(err)
	}
}

func (ctx *Context) Delete(ts string) {
	ctx.bot.client.DeleteMessage(ctx.Channel, ts)
}

func (ctx *Context) Upload(params slack.FileUploadParameters) error {
	_, err := ctx.bot.client.UploadFile(params)
	return err
}

func (mb MessageButton) ToAction() slack.AttachmentAction {
	return slack.AttachmentAction{
		Name:  mb.Name,
		Text:  mb.Text,
		Type:  "button",
		Value: mb.Value,
	}
}

func (mm MessageMenu) ToAction() slack.AttachmentAction {
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
