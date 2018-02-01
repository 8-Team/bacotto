package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/8-team/bacotto/conf"
	"github.com/8-team/bacotto/db"
	"github.com/8-team/bacotto/erp"
	"github.com/Sirupsen/logrus"
)

var log = logrus.WithField("app", "api")

func isAuthorize(r *http.Request) bool {
	serial := r.FormValue("serial")
	otp := r.FormValue("otp")

	if _, err := Authorize(serial, otp); err != nil {
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

func prjList(w http.ResponseWriter, r *http.Request) {
	if !isAuthorize(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var err error
	var js []byte
	prefered := r.FormValue("prefered")
	if strings.Compare(prefered, "true") == 0 {
		// read from db
		device := []db.Project{}
		if db.DB.Find(&device).RecordNotFound() {
			log.Errorln("Unable to get connection to ERP")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		js, err = json.Marshal(device)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		var prjs []erp.Project
		err = erp.Open()
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
		js, err = json.Marshal(prjs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

type OttoProjectEntry struct {
	ProjectID uint
	StartDate time.Time
	Duration  time.Duration
	Serial    string
	OTP       string
}

func registerEntry(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t OttoProjectEntry
	err := decoder.Decode(&t)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	log.Info(t.ProjectID, t.StartDate, t.Duration)

	otto, err := Authorize(t.Serial, t.OTP)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if otto.UserID == nil {
		log.Error("No user registered for this Otto.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// search project by id and add/update entry
	prj := db.Project{}
	if db.DB.Find(&prj, t.ProjectID).RecordNotFound() {
		log.Error("Invalid Project ID: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info(prj)

	log.Info("User id", otto.UserID, otto.Serial)
	entry := db.ProjectEntry{UserID: 0, ProjectID: t.ProjectID,
		StartTime: t.StartDate,
		EndTime:   t.StartDate.Add(t.Duration)}

	if db.DB.NewRecord(rec) {
		db.DB.Create(&rec)
	}
	if err := db.DB.Create(&entry, "project_id = ?", entry.ProjectID).Error; err != nil {
		log.Error("Unable to add entry", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ListenAndServe(addr string) error {
	http.HandleFunc("/ping", pong)
	http.HandleFunc("/prjlist", prjList)
	http.HandleFunc("/register", registerEntry)

	if conf.UseHTTPS() {
		return http.ListenAndServeTLS(addr,
			conf.GetCertFilePath(),
			conf.GetKeyfilePath(), nil)
	}

	return http.ListenAndServe(addr, nil)
}
