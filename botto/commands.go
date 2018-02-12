package botto

import (
	"strings"

	"github.com/nlopes/slack"
)

func (ctx *context) showHelp(ev *slack.MessageEvent) {
	text := `Here's some stuff you can do:
` + "`botto help` to show this help message" + `
` + "`botto add project` to start managing a project using your Otto" + `
` + "`botto list projects` to show a list of your managed projects" + `
` + "`botto report` to show a report of today's activities"

	ctx.Send(text)
}

func (ctx *context) parseCommand(ev *slack.MessageEvent) {
	cmd := strings.TrimLeft(strings.TrimSpace(ev.Text), "botto ")

	switch cmd {
	case "add project":
		ctx.pickProject(ev)
	case "list projects":
		ctx.listProjects(ev)
	case "report":
		ctx.showReport(ev)
	default:
		ctx.showHelp(ev)
	}
}
