package botto

import (
	"regexp"

	"github.com/8-team/bacotto/db"
	"github.com/nlopes/slack"
	"github.com/plorefice/slackbot"
	"github.com/pquerna/otp/totp"
)

type registrationContext struct {
	device *db.Otto
	otp    string
}

func greetUser(bot *slackbot.Bot, msg *slack.Msg, ctx interface{}) bool {
	text := `Hi! Welcome to Otto! I'm Botto and I'll be your guide.
I see this is the first time using Otto, so I'll get you up and running.
To get started, please input your Otto's serial number. You can find it on the back of your device.`

	bot.Message(msg.Channel, text)
	return true
}

func inputSerial(bot *slackbot.Bot, msg *slack.Msg, ctx interface{}) bool {
	reg := ctx.(*registrationContext)
	reg.device = new(db.Otto)

	re := regexp.MustCompile("^[a-fA-F0-9]{6}$")
	if !re.MatchString(msg.Text) {
		bot.Message(msg.Channel, "The serial number must be a 6-digit hex number.")
		return false
	}

	if db.DB.First(reg.device, "serial = ?", msg.Text).RecordNotFound() {
		bot.Message(msg.Channel, "Sorry, I can't find a matching serial. Is there a typo somewhere?")
		return false
	}

	bot.Message(msg.Channel, "You are doing great! Now input the OTP code from your Otto.")
	return true
}

func inputOtp(bot *slackbot.Bot, msg *slack.Msg, ctx interface{}) bool {
	reg := ctx.(*registrationContext)

	re := regexp.MustCompile("^\\d{6}$")
	if !re.MatchString(msg.Text) {
		bot.Message(msg.Channel, "The OTP must be a 6-digit number.")
		return false
	}

	if !totp.Validate(msg.Text, reg.device.OTPSecret) {
		bot.Message(msg.Channel, "Sorry, this OTP is not valid, try again.")
		return false
	}

	user := &db.User{
		Username: msg.User,
	}

	if db.DB.Create(user).Error != nil {
		bot.Message(msg.Channel, "There was an error during registration, try again.")
		return false
	}

	reg.device.UserID = &user.ID
	db.DB.Save(reg.device)

	bot.Message(msg.Channel, "You are good to go, thank you for using Otto!")
	return true
}

func onRegistrationRequest(bot *slackbot.Bot, msg *slack.Msg) bool {
	var user db.User
	return db.DB.First(&user, "username = ?", msg.User).RecordNotFound()
}

var registrationFlow *slackbot.Flow

func init() {
	greetState := slackbot.NewState("greet_user", greetUser).To("input_serial").Build()
	serialState := slackbot.NewState("input_serial", inputSerial).To("input_otp").Build()
	otpState := slackbot.NewState("input_otp", inputOtp).Build()

	registrationFlow = slackbot.NewFlowWithContext("registration_flow", &registrationContext{}).
		AddStates(greetState, serialState, otpState).
		FilterBy(slackbot.DMFilter).
		SetGuard(onRegistrationRequest).
		Build("greet_user")
}
