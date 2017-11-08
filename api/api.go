package api

import (
	"net/http"

	"github.com/8-team/bacotto/conf"
	"github.com/8-team/bacotto/db"
	"github.com/Sirupsen/logrus"
	"github.com/pquerna/otp/totp"
)

var log = logrus.WithField("app", "api")

func pong(w http.ResponseWriter, r *http.Request) {
	var otto db.Otto

	serial := r.FormValue("serial")
	otp := r.FormValue("otp")

	if err := db.DB.First(&otto, "serial = ?", serial).Error; err != nil {
		log.Errorf("http request from invalid Otto serial %s: %s", serial, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if totp.Validate(otp, otto.OTPSecret) {
		w.WriteHeader(http.StatusOK)
	} else {
		log.Warningf("invalid OTP code %s from %s", otp, serial)
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func ListenAndServe(addr string) error {
	http.HandleFunc("/ping", pong)

	return http.ListenAndServeTLS(addr,
		conf.GetCertFilePath(),
		conf.GetKeyfilePath(), nil)
}
