package botto

import (
	"errors"
	"regexp"
	"strings"

	"github.com/8-team/bacotto/db"
)

func (uc *userContext) validateSerial(serial string) error {
	re := regexp.MustCompile("^[a-fA-F0-9]{6}$")

	if !re.MatchString(serial) {
		return errors.New("Sorry, the serial number must be a 6-digit hex number")
	}

	uc.currentDevice = new(db.Otto)
	if db.DB.First(uc.currentDevice, "serial = ?", serial).RecordNotFound() {
		return errors.New("Sorry, I can't find a matching serial. Could you double-check it?")
	}

	return nil
}

func (uc *userContext) validateOtp(otp string) error {
	re := regexp.MustCompile("^\\d{6}$")

	if !re.MatchString(otp) {
		return errors.New("Sorry, the OTP must be a 6-digit number")
	}

	if _, err := db.Authorize(uc.currentDevice.Serial, otp); err != nil {
		return errors.New("Sorry, this OTP is not valid, try again")
	}

	return nil
}

func (uc *userContext) createUser(username string) error {
	user := &db.User{
		Username: username,
	}

	if db.DB.Create(user).Error != nil {
		return errors.New("There was an error during registration, try again")
	}

	uc.currentDevice.UserID = &user.ID
	db.DB.Save(uc.currentDevice)

	return nil
}

func (uc *userContext) registerUser(bot *slackbot, ev contextEvent) {
	text := `Hi! I'm Botto and I'll be your guide in the magic world of Otto :8ball:
Looks like this is your first time using your Otto, so let's get you up and running!
To get started, please input your Otto's serial number. You can find it on the back of your device.`

	bot.Message(ev.channel(), text)
	uc.dispatcher = uc.inputSerial
}

func (uc *userContext) inputSerial(bot *slackbot, ev contextEvent) {
	if err := uc.validateSerial(strings.TrimSpace(ev.text())); err != nil {
		bot.Message(ev.channel(), err.Error())
	} else {
		bot.Message(ev.channel(), "You are doing great! Now input the OTP code from your Otto.")
		uc.dispatcher = uc.inputOtp
	}
}

func (uc *userContext) inputOtp(bot *slackbot, ev contextEvent) {
	if err := uc.validateOtp(strings.TrimSpace(ev.text())); err != nil {
		bot.Message(ev.channel(), err.Error())
		return
	}

	if err := uc.createUser(ev.user()); err != nil {
		bot.Message(ev.channel(), err.Error())
		return
	}

	bot.Message(ev.channel(), "Fantastic, let's try adding a project to your Otto!")

	uc.dispatcher = uc.pickProject
	uc.dispatcher(bot, ev)
}
