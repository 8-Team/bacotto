package botto

/*
func helpUser(bot *slackbot.Bot, msg *slack.Msg, ctx interface{}) bool {
	helpText := `Here's some stuff you can do:
` + "`botto help` to show this help message" + `
` + "`botto 42` for an easter egg"

	bot.Message(msg.Channel, helpText)
	return true
}

func unknownCmd(bot *slackbot.Bot, msg *slack.Msg, ctx interface{}) bool {
	bot.Message(msg.Channel, "Sorry, I don't know how to do that.")
	return true
}

func onHelpRequest(bot *slackbot.Bot, msg *slack.Msg) bool {
	re := regexp.MustCompile("^(b?otto )?help$")
	return re.MatchString(strings.TrimSpace(msg.Text))
}

var helpFlow *slackbot.Flow
var unknownCommandFlow *slackbot.Flow

func init() {
	helpFlow = slackbot.NewFlow("help_flow").
		AddStates(slackbot.NewState("help_user", helpUser).Build()).
		FilterBy(slackbot.DMFilter).
		SetGuard(onHelpRequest).
		Build("help_user")

	unknownCommandFlow = slackbot.NewFlow("unk_cmd_flow").
		AddStates(slackbot.NewState("unk_cmd", unknownCmd).Build()).
		FilterBy(slackbot.DMFilter).
		SetGuard(func(bot *slackbot.Bot, msg *slack.Msg) bool { return true }).
		Build("unk_cmd")
}
*/
