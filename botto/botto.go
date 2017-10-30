package botto

import (
	"errors"

	"github.com/plorefice/slackbot"
)

// ListenAndServe starts the bot using the given API token
func ListenAndServe(token string) error {
	bot := slackbot.New(token, slackbot.Config{})
	if bot == nil {
		return errors.New("Could not create Slack bot")
	}

	return bot.Start()
}
