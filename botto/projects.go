package botto

import (
	"fmt"

	"github.com/8-team/bacotto/db"
	"github.com/nlopes/slack"
)

func (ctx *context) pickProject(ev *slack.MessageEvent) {
	var projects []db.Project

	if err := db.DB.Find(&projects).Error; err != nil {
		ctx.Send("There was a problem retrieving your projects, try later.")
		return
	}

	menu := messageMenu{
		Name:   "projects",
		Values: make(map[string]string),
	}

	for _, p := range projects {
		menu.Values[p.Name] = p.Name
	}

	msg := interactiveMessage{
		Elements: []interactiveElement{menu},
	}

	ctx.SendInteractive("Here is a list of your recent projects, "+
		"select the ones you want to see on your device:", msg,
		func(resp *slack.AttachmentActionCallback) {
			project := new(db.Project)
			name := resp.Actions[0].SelectedOptions[0].Value

			if err := db.DB.First(project, "name = ?", name).Error; err != nil {
				ctx.Send("Invalid project selected")
				return
			}

			ctx.user.Projects = append(ctx.user.Projects, project)
			db.DB.Save(ctx.user)

			ctx.Update(resp, interactiveMessage{
				Text: ":heavy_check_mark: Project successfully added!",
			})
		})
}

func (ctx *context) listProjects(ev *slack.MessageEvent) {
	projects, err := db.GetUserProjects(ctx.user)
	if err != nil {
		ctx.Send("There was a problem retrieving your projects, try later.")
		return
	}

	resp := "Here's a list of your currently tracked projects:\n"
	for _, prj := range projects {
		resp += fmt.Sprintf("* %s\n", prj.Name)
	}

	ctx.Send(resp)
}
