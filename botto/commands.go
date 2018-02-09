package botto

import (
	"strings"
)

func (uc *userContext) showHelp(bot *slackbot, ev contextEvent) {
	text := `Here's some stuff you can do:
` + "`botto help` to show this help message" + `
` + "`botto add project` to start managing a project using your Otto" + `
` + "`botto list projects` to show a list of your managed projects" + `
` + "`botto report` to show a report of today's activities"

	bot.Message(ev.channel(), text)
}

func (uc *userContext) parseCommand(bot *slackbot, ev contextEvent) {
	cmd := strings.TrimLeft(strings.TrimSpace(ev.text()), "botto ")

	switch cmd {
	case "add project":
		uc.pickProject(bot, ev)
	case "list projects":
		uc.listProjects(bot, ev)
	case "report":
		uc.showReport(bot, ev)
	default:
		uc.showHelp(bot, ev)
	}
}
