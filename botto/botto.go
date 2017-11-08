package botto

import (
	"github.com/Sirupsen/logrus"
	"github.com/plorefice/slackbot"
)

var log = logrus.WithField("app", "botto")

// ListenAndServe starts the bot using the given API token
func ListenAndServe(token string) error {
	bot, err := slackbot.New(token, slackbot.Config{})
	if err != nil {
		return err
	}

	// Register higher priority flows first
	bot.RegisterFlow(registrationFlow)
	bot.RegisterFlow(helpFlow)
	bot.RegisterFlow(unknownCommandFlow)

	return bot.Start()
}
