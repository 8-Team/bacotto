package api

import (
	"fmt"

	"github.com/8-team/bacotto/conf"
	"github.com/8-team/bacotto/db"
	"github.com/pquerna/otp/totp"
)

// Authorize authenticates a request from an Otto by using the provided OTP code.
func Authorize(serial string, otp string) (db.Otto, error) {
	var otto db.Otto

	if err := db.DB.First(&otto, "serial = ?", serial).Error; err != nil {
		return otto, fmt.Errorf("http request from invalid Otto serial %s: %s", serial, err)
	}

	if !conf.VerifyMsg() {
		log.Warn("You are in debug mode, no serial and OPT validation! check config file.")
		return otto, nil
	}

	if !totp.Validate(otp, otto.OTPSecret) {
		return otto, fmt.Errorf("invalid OTP code %s from %s", otp, serial)
	}

	return otto, nil
}
