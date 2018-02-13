package botto

import (
	"bytes"
	"image/png"
	"time"

	"github.com/8-team/bacotto/db"
	"github.com/8-team/bacotto/graph"
	"github.com/8-team/bacotto/slackbot"
	"github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
)

func ShowReport(ctx *slackbot.Context, ev *slack.MessageEvent) error {
	badJuju := "Sorry, there was a problem with your report, try later"
	today := time.Now()
	user := ctx.Data.(*db.User)

	entries, err := db.GetUserEntries(user, today, today)
	if err != nil {
		ctx.Send(badJuju)
		return err
	}

	pc := graph.NewPunchcard()
	pc.SetDay(today)

	/* TODO: make this generic */
	loc, _ := time.LoadLocation("Europe/Rome")

	for _, e := range entries {
		var prj db.Project

		if err := db.DB.First(&prj, e.ProjectID).Error; err != nil {
			logrus.Warnf("error retrieving project: %s", err)
			continue
		}

		pc.AddTask(prj.Name, e.StartTime.In(loc), e.EndTime.In(loc))
	}

	img := pc.Rasterize()
	buf := bytes.NewBuffer(make([]byte, 0))
	png.Encode(buf, img)

	params := slack.FileUploadParameters{
		Title:    "Your daily report",
		Channels: []string{ctx.Channel},
		Filetype: "png",
		Filename: "report.png",
		Reader:   buf,
	}

	ts := ctx.Send("I'm working on it, just a sec")

	if err := ctx.Upload(params); err != nil {
		ctx.Send(badJuju)
	}

	ctx.Delete(ts)

	return nil
}
