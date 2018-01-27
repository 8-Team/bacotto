package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/8-team/bacotto/conf"
	"github.com/8-team/bacotto/erp"
	"github.com/Sirupsen/logrus"
)

var log = logrus.WithField("app", "api")

func isAuthorize(r *http.Request) bool {
	serial := r.FormValue("serial")
	otp := r.FormValue("otp")

	if err := Authorize(serial, otp); err != nil {
		log.Errorln(err)
		return false
	}
	return true
}

func pong(w http.ResponseWriter, r *http.Request) {
	if !isAuthorize(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func prjlist(w http.ResponseWriter, r *http.Request) {
	if !isAuthorize(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	prefered := r.FormValue("prefered")

	var prjs []erp.Project
	if strings.Compare(prefered, "true") == 0 {
		// read from db
		log.Info("No db")
	} else {
		err := erp.Open()
		if err != nil {
			log.Errorln("Unable to get connection to ERP")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		prjs, err = erp.ListProjects()
		if err != nil {
			log.Errorln("Unable to get prj list from ERP")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	js, err := json.Marshal(prjs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func ListenAndServe(addr string) error {
	http.HandleFunc("/ping", pong)
	http.HandleFunc("/prjlist", prjlist)

	if conf.UseHTTPS() {
		return http.ListenAndServeTLS(addr,
			conf.GetCertFilePath(),
			conf.GetKeyfilePath(), nil)
	}

	return http.ListenAndServe(addr, nil)
}
