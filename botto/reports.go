package botto

import (
	"bytes"
	"image/png"
	"time"

	"github.com/8-team/bacotto/db"
	"github.com/8-team/bacotto/graph"
	"github.com/nlopes/slack"
)

func (ctx *context) showReport(ev *slack.MessageEvent) {
	badJuju := "Sorry, there was a problem with your report, try later"
	today := time.Now()

	entries, err := db.GetUserEntries(ctx.user, today, today)
	if err != nil {
		ctx.Send(badJuju)
		return
	}

	pc := graph.NewPunchcard()
	pc.SetDay(today)

	for _, e := range entries {
		var prj db.Project

		if err := db.DB.First(&prj, e.ProjectID).Error; err != nil {
			log.Warnf("error retrieving project: %s", err)
			continue
		}

		pc.AddTask(prj.Name, e.StartTime, e.EndTime)
	}

	img := pc.Rasterize()
	buf := bytes.NewBuffer(make([]byte, 0))
	png.Encode(buf, img)

	params := slack.FileUploadParameters{
		Title:    "Your daily report",
		Channels: []string{ctx.channel},
		Filetype: "png",
		Filename: "report.png",
		Reader:   buf,
	}

	ts := ctx.Send("I'm working on it, just a sec")

	if _, err := bot.client.UploadFile(params); err != nil {
		ctx.Send(badJuju)
	}

	ctx.bot.client.DeleteMessage(ev.Channel, ts)
}
