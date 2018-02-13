package botto

import (
	"errors"
	"regexp"
	"strings"

	"github.com/8-team/bacotto/db"
	"github.com/8-team/bacotto/slackbot"
	"github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
)

func validateSerial(user *db.User, serial string) error {
	re := regexp.MustCompile("^[a-fA-F0-9]{6}$")

	if !re.MatchString(serial) {
		return errors.New("Sorry, the serial number must be a 6-digit hex number")
	}

	if db.DB.First(&user.Otto, "serial = ?", serial).RecordNotFound() {
		return errors.New("Sorry, I can't find a matching serial. Could you double-check it?")
	}

	return nil
}

func validateOtp(serial, otp string) error {
	re := regexp.MustCompile("^\\d{6}$")

	if !re.MatchString(otp) {
		return errors.New("Sorry, the OTP must be a 6-digit number")
	}

	if _, err := db.Authorize(serial, otp); err != nil {
		return errors.New("Sorry, this OTP is not valid, try again")
	}

	return nil
}

func registerUser(ctx *slackbot.Context) {
	user := ctx.Data.(*db.User)

	ctx.Send(`Hi! I'm Botto and I'll be your guide in the magic world of Otto :8ball:
Looks like this is your first time using your Otto, so let's get you up and running!
To get started, please input your Otto's serial number. You can find it on the back of your device.`)

	for {
		msg := ctx.Receive()
		serial := strings.TrimSpace(msg.Text)

		if err := validateSerial(user, serial); err != nil {
			ctx.Send(err.Error())
		} else {
			break
		}
	}

	ctx.Send("You are doing great! Now input the OTP code from your Otto.")

	for {
		msg := ctx.Receive()
		otp := strings.TrimSpace(msg.Text)

		if err := validateOtp(user.Otto.Serial, otp); err != nil {
			ctx.Send(err.Error())
		} else {
			break
		}
	}

	if err := db.DB.Save(user).Error; err != nil {
		ctx.Send(err.Error())
		return
	}

	ctx.Send("Fantastic, let's try adding a project to your Otto!")

	PickProject(ctx, nil)

	ctx.Send("You are set and ready to go! Type `help` to show a list of possible commands.")
}

func CheckUserPresence(ctx *slackbot.Context, ev *slack.MessageEvent) error {
	ctx.Data = new(db.User)
	user := ctx.Data.(*db.User)

	if err := db.DB.Preload("Otto").First(user, "username = ?", ev.Username).Error; err != nil {
		logrus.Debugln("User not found in DB, proceeding with registration")
		registerUser(ctx)
	}
	return nil
}
