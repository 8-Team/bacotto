package botto

import (
	"github.com/plorefice/slackbot"
)

// ListenAndServe starts the bot using the given API token
func ListenAndServe(token string) error {
	bot, err := slackbot.New(token, slackbot.Config{})
	if err != nil {
		return err
	}

	bot.RegisterFlow(registrationFlow)

	return bot.Start()
}
