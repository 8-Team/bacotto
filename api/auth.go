package api

import (
	"fmt"

	"github.com/8-team/bacotto/conf"
	"github.com/8-team/bacotto/db"
	"github.com/pquerna/otp/totp"
)

// Authorize authenticates a request from an Otto by using the provided OTP code.
func Authorize(serial string, otp string) error {
	var otto db.Otto

	if !conf.VerifyMsg() {
		log.Warn("You are in debug mode, no serial and OPT validation! check config file.")
		return nil
	}

	if err := db.DB.First(&otto, "serial = ?", serial).Error; err != nil {
		return fmt.Errorf("http request from invalid Otto serial %s: %s", serial, err)
	}

	if !totp.Validate(otp, otto.OTPSecret) {
		return fmt.Errorf("invalid OTP code %s from %s", otp, serial)
	}
	return nil
}
