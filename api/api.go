package api

import (
	"net/http"

	"github.com/8-team/bacotto/conf"
	"github.com/Sirupsen/logrus"
)

var log = logrus.WithField("app", "api")

func pong(w http.ResponseWriter, r *http.Request) {
	serial := r.FormValue("serial")
	otp := r.FormValue("otp")

	if err := Authorize(serial, otp); err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func ListenAndServe(addr string) error {
	http.HandleFunc("/ping", pong)

	if conf.UseHTTPS() {
		return http.ListenAndServeTLS(addr,
			conf.GetCertFilePath(),
			conf.GetKeyfilePath(), nil)
	}

	return http.ListenAndServe(addr, nil)
}
