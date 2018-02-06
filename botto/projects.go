package botto

import (
	"fmt"

	"github.com/8-team/bacotto/db"
)

func (uc *userContext) projectSelected(bot *slackbot, resp *interactiveResponse) {
	project := new(db.Project)
	name := resp.Actions[0].SelectedOptions[0].Value

	if err := db.DB.First(project, "name = ?", name).Error; err != nil {
		bot.Message(resp.channel(), "Invalid project selected")
		return
	}

	uc.user.Projects = append(uc.user.Projects, project)
	db.DB.Save(uc.user)

	bot.Update(resp, interactiveMessage{
		Text: ":heavy_check_mark: Project successfully added!",
	})
}

func (uc *userContext) pickProject(bot *slackbot, ev contextEvent) {
	var projects []db.Project

	if err := db.DB.Find(&projects).Error; err != nil {
		bot.Message(ev.channel(), "There was a problem retrieving your projects, try later.")
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
		Callback: uc.projectSelected,
		Elements: []interactiveElement{menu},
	}

	bot.InteractiveMessage(ev.channel(), "Here is a list of your recent projects, "+
		"select the ones you want to see on your device:", msg)

	uc.dispatcher = uc.parseCommand
}

func (uc *userContext) listProjects(bot *slackbot, ev contextEvent) {
	projects, err := db.GetProjects(uc.user.Otto.Serial)
	if err != nil {
		bot.Message(ev.channel(), "There was a problem retrieving your projects, try later.")
		return
	}

	resp := "Here's a list of your currently tracked projects:\n"
	for _, prj := range projects {
		resp += fmt.Sprintf("* %s\n", prj.Name)
	}

	bot.Message(ev.channel(), resp)

	uc.dispatcher = uc.parseCommand
}
