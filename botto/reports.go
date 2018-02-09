package botto

import (
	"bytes"
	"image/png"
	"time"

	"github.com/8-team/bacotto/graph"
	"github.com/nlopes/slack"
)

func (uc *userContext) showReport(bot *slackbot, ev contextEvent) {
	pc := graph.NewPunchcard()
	pc.SetDay(time.Now())

	pc.AddTask("Mock entry", time.Date(0, 0, 0, 9, 42, 0, 0, time.Local),
		time.Date(0, 0, 0, 13, 2, 0, 0, time.Local))

	img := pc.Rasterize()
	buf := bytes.NewBuffer(make([]byte, 0))
	png.Encode(buf, img)

	params := slack.FileUploadParameters{
		Title:    "Your daily report",
		Channels: []string{ev.channel()},
		Filetype: "png",
		Filename: "report.png",
		Reader:   buf,
	}

	ts := bot.Message(ev.channel(), "I'm working on it, just a sec")

	if _, err := bot.client.UploadFile(params); err != nil {
		bot.Message(ev.channel(), "Sorry, there was a problem with your report, try later")
	}

	bot.client.DeleteMessage(ev.channel(), ts)
}
