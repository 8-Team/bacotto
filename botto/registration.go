package botto

import (
	"errors"
	"regexp"
	"strings"

	"github.com/8-team/bacotto/db"
	"github.com/nlopes/slack"
)

func (ctx *context) validateSerial(serial string) error {
	re := regexp.MustCompile("^[a-fA-F0-9]{6}$")

	if !re.MatchString(serial) {
		return errors.New("Sorry, the serial number must be a 6-digit hex number")
	}

	if db.DB.First(&ctx.user.Otto, "serial = ?", serial).RecordNotFound() {
		return errors.New("Sorry, I can't find a matching serial. Could you double-check it?")
	}

	return nil
}

func (ctx *context) validateOtp(otp string) error {
	re := regexp.MustCompile("^\\d{6}$")

	if !re.MatchString(otp) {
		return errors.New("Sorry, the OTP must be a 6-digit number")
	}

	if _, err := db.Authorize(ctx.user.Otto.Serial, otp); err != nil {
		return errors.New("Sorry, this OTP is not valid, try again")
	}

	return nil
}

func (ctx *context) registerUser(ev *slack.MessageEvent) {
	ctx.Send(`Hi! I'm Botto and I'll be your guide in the magic world of Otto :8ball:
Looks like this is your first time using your Otto, so let's get you up and running!
To get started, please input your Otto's serial number. You can find it on the back of your device.`)

	ev = ctx.Wait()
	serial := strings.TrimSpace(ev.Text)

	for err := ctx.validateSerial(serial); err != nil; {
		ctx.Send(err.Error())
	}

	ctx.Send("You are doing great! Now input the OTP code from your Otto.")

	ev = ctx.Wait()
	otp := strings.TrimSpace(ev.Text)

	for err := ctx.validateOtp(otp); err != nil; {
		ctx.Send(err.Error())
	}

	if err := db.DB.Save(ctx.user).Error; err != nil {
		ctx.Send(err.Error())
		return
	}

	ctx.Send("Fantastic, let's try adding a project to your Otto!")

	ctx.pickProject(nil)

	ctx.Send("You are set and ready to go! Type `help` to show a list of possible commands.")
}
