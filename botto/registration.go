package botto

import (
	"errors"
	"regexp"
	"strings"

	"github.com/8-team/bacotto/db"
	"github.com/nlopes/slack"
	"github.com/plorefice/slackbot"
	"github.com/pquerna/otp/totp"
)

type registrationContext struct {
	device *db.Otto
	otp    string
}

func (rc *registrationContext) validateSerial(serial string) error {
	re := regexp.MustCompile("^[a-fA-F0-9]{6}$")

	if !re.MatchString(serial) {
		return errors.New("Sorry, the serial number must be a 6-digit hex number")
	}

	rc.device = new(db.Otto)
	if db.DB.First(rc.device, "serial = ?", serial).RecordNotFound() {
		return errors.New("Sorry, I can't find a matching serial. Could you double-check it?")
	}

	return nil
}

func (rc *registrationContext) validateOtp(otp string) error {
	re := regexp.MustCompile("^\\d{6}$")

	if !re.MatchString(otp) {
		return errors.New("Sorry, the OTP must be a 6-digit number")
	}

	if !totp.Validate(otp, rc.device.OTPSecret) {
		return errors.New("Sorry, this OTP is not valid, try again")
	}

	return nil
}

func (rc *registrationContext) createUser(username string) error {
	user := &db.User{
		Username: username,
	}

	if db.DB.Create(user).Error != nil {
		return errors.New("There was an error during registration, try again")
	}

	rc.device.UserID = &user.ID
	db.DB.Save(rc.device)

	return nil
}

func greetUser(bot *slackbot.Bot, msg *slack.Msg, ctx interface{}) bool {
	text := `Hi! I'm Botto and I'll be your guide in the magic world of Otto :8ball:
I see this is the first time using your Otto, so let's get you up and running!
To get started, please input your Otto's serial number. You can find it on the back of your device.`

	bot.Message(msg.Channel, text)
	return true
}

func inputSerial(bot *slackbot.Bot, msg *slack.Msg, ctx interface{}) bool {
	reg := ctx.(*registrationContext)

	if err := reg.validateSerial(strings.TrimSpace(msg.Text)); err != nil {
		bot.Message(msg.Channel, err.Error())
		return false
	}

	bot.Message(msg.Channel, "You are doing great! Now input the OTP code from your Otto.")
	return true
}

func inputOtp(bot *slackbot.Bot, msg *slack.Msg, ctx interface{}) bool {
	reg := ctx.(*registrationContext)

	if err := reg.validateOtp(strings.TrimSpace(msg.Text)); err != nil {
		bot.Message(msg.Channel, err.Error())
		return false
	}

	if err := reg.createUser(msg.User); err != nil {
		bot.Message(msg.Channel, err.Error())
		return false
	}

	bot.Message(msg.Channel, "Fantastic, let's move on!")
	return true
}

/*
	# Example code to show an interactive project list:

	prjs, err := erp.ListProjects()
	if err != nil {
		bot.Message("channel string", msg string)
	}

	menu := slackbot.MessageMenu{
		Name:   "projects",
		Values: make(map[string]string),
	}
	for _, p := range prjs {
		menu.Values[p.Name] = p.Name
	}

	fmt := slackbot.MessageFormat{
		Callback: "project_selection",
		Elements: []slackbot.InteractiveElement{
			menu,
			slackbot.MessageButton{
				Name:  "confirm",
				Text:  "I'm done",
				Value: "confirm",
			},
		},
	}

	bot.InteractiveMessage(msg.Channel, "Here is a list of your recent projects, "+
		"select the ones you want to see on your device", fmt)
*/

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
