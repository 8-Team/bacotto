package botto

import (
	"github.com/8-team/bacotto/erp"
)

func (uc *userContext) pickProjects(bot *slackbot, ev contextEvent) {
	prjs, err := erp.ListProjects()
	if err != nil {
		bot.Message(ev.channel(), "There was a problem retrieving your projects, try later.")
		uc.dispatcher = uc.genericDispatcher
		return
	}

	menu := messageMenu{
		Name:   "projects",
		Values: make(map[string]string),
	}
	for _, p := range prjs {
		menu.Values[p.Name] = p.Name
	}

	fmt := messageFormat{
		Callback: "project_selection",
		Elements: []interactiveElement{
			menu,
			messageButton{
				Name:  "confirm",
				Text:  "I'm done",
				Value: "confirm",
			},
		},
	}

	bot.InteractiveMessage(ev.channel(), "Here is a list of your recent projects, "+
		"select the ones you want to see on your device:", fmt)
}
