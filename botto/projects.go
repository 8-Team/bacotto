package botto

import (
	"fmt"

	"github.com/8-team/bacotto/db"
	"github.com/8-team/bacotto/slackbot"
	"github.com/nlopes/slack"
)

func PickProject(ctx *slackbot.Context, ev *slack.MessageEvent) error {
	var projects []db.Project

	user := ctx.Data.(*db.User)

	if err := db.DB.Find(&projects).Error; err != nil {
		ctx.Send("There was a problem retrieving your projects, try later.")
		return err
	}

	menu := slackbot.MessageMenu{
		Name:   "projects",
		Values: make(map[string]string),
	}

	for _, p := range projects {
		menu.Values[p.Name] = p.Name
	}

	msg := slackbot.InteractiveMessage{
		Elements: []slackbot.InteractiveElement{menu},
	}

	text := "Here is a list of your recent projects, " +
		"select the ones you want to see on your device:"

	ctx.SendInteractive(text, msg, func(resp *slack.AttachmentActionCallback) {
		project := new(db.Project)
		name := resp.Actions[0].SelectedOptions[0].Value

		if err := db.DB.First(project, "name = ?", name).Error; err != nil {
			ctx.Send("Invalid project selected")
			return
		}

		user.Projects = append(user.Projects, project)
		db.DB.Save(user)

		ctx.Update(resp, slackbot.InteractiveMessage{
			Text: ":heavy_check_mark: Project successfully added!",
		})
	})

	return nil
}

func ListProjects(ctx *slackbot.Context, ev *slack.MessageEvent) error {
	user := ctx.Data.(*db.User)

	projects, err := db.GetUserProjects(user)
	if err != nil {
		ctx.Send("There was a problem retrieving your projects, try later.")
		return err
	}

	resp := "Here's a list of your currently tracked projects:\n"
	for _, prj := range projects {
		resp += fmt.Sprintf("* %s\n", prj.Name)
	}

	ctx.Send(resp)

	return nil
}
