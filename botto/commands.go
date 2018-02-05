package botto

import (
	"strings"
)

func (uc *userContext) showHelp(bot *slackbot, ev contextEvent) {
	text := `Here's some stuff you can do:
` + "`botto help` to show this help message" + `
` + "`botto add project` to manage a project using your Otto"

	bot.Message(ev.channel(), text)
}

func (uc *userContext) parseCommand(bot *slackbot, ev contextEvent) {
	cmd := strings.TrimLeft(strings.TrimSpace(ev.text()), "botto ")

	switch cmd {
	case "add project":
		uc.pickProject(bot, ev)
		break

	default:
		uc.showHelp(bot, ev)
		break
	}
}
