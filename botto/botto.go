package botto

import (
	"fmt"

	"github.com/nlopes/slack"
	"github.com/plorefice/slackbot"
)

type testContext struct {
	a string
}

func testGuard(bot *slackbot.Bot, msg *slack.Msg) bool {
	fmt.Println("guard invoked for", msg.Text)
	return true
}

func state1_act(bot *slackbot.Bot, msg *slack.Msg, ctx interface{}) bool {
	tc := ctx.(*testContext)

	bot.Message(msg.Channel, "You wrote in state1: "+msg.Text+" with context = "+tc.a)
	tc.a = "ciao"

	return true
}

func state2_act(bot *slackbot.Bot, msg *slack.Msg, ctx interface{}) bool {
	tc := ctx.(*testContext)

	bot.Message(msg.Channel, "You wrote in state2: "+msg.Text+" with context = "+tc.a)
	tc.a = "hola"

	return true
}

func state3_act(bot *slackbot.Bot, msg *slack.Msg, ctx interface{}) bool {
	tc := ctx.(*testContext)

	bot.Message(msg.Channel, "You wrote in state3 (exiting): "+msg.Text+" with context = "+tc.a)
	tc.a = "hello"

	return true
}

// ListenAndServe starts the bot using the given API token
func ListenAndServe(token string) error {
	bot, err := slackbot.New(token, slackbot.Config{})
	if err != nil {
		return err
	}

	state1 := slackbot.NewState("state_1", state1_act).To("state_2").Build()
	state2 := slackbot.NewState("state_2", state2_act).To("state_3").Build()
	state3 := slackbot.NewState("state_3", state3_act).Build()

	testFlow := slackbot.NewFlowWithContext("test_flow", &testContext{}).
		AddStates(state1, state2, state3).
		FilterBy(slackbot.DMFilter).
		SetGuard(testGuard).
		Build("state_1")

	bot.RegisterFlow(testFlow)

	return bot.Start()
}
