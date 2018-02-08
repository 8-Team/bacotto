package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/8-team/bacotto/botto"
	"github.com/8-team/bacotto/conf"
	"github.com/8-team/bacotto/db"
	"github.com/Sirupsen/logrus"
)

var log = logrus.WithField("app", "api")

func isAuthorized(r *http.Request) bool {
	serial := r.FormValue("serial")
	otp := r.FormValue("otp")

	if _, err := db.Authorize(serial, otp); err != nil {
		log.Errorln(err)
		return false
	}
	return true
}

func pong(w http.ResponseWriter, r *http.Request) {
	if !isAuthorized(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func listProjects(w http.ResponseWriter, r *http.Request) {
	serial := r.FormValue("serial")

	if !isAuthorized(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	projects, err := db.GetProjects(serial)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if projects == nil {
		projects = []*db.Project{}
	}

	js, err := json.Marshal(projects)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// OttoProjectEntry ... EtryProject from otto
type OttoProjectEntry struct {
	ProjectID uint
	StartDate time.Time
	Duration  time.Duration
	Serial    string
	OTP       string
}

func registerEntry(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var items []OttoProjectEntry
	err := decoder.Decode(&items)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	for _, t := range items {
		log.Info(t)
		otto, err := db.Authorize(t.Serial, t.OTP)
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

		log.Info("User id", otto.UserID, otto.Serial)
		entry := db.ProjectEntry{UserID: *otto.UserID, ProjectID: t.ProjectID,
			StartTime: t.StartDate,
			EndTime:   t.StartDate.Add(t.Duration * time.Minute)}

		if db.DB.NewRecord(entry) {
			if err := db.DB.Create(&entry).Error; err != nil {
				log.Error("Unable to add entry", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}

}

// OttoStatsProject ... Registered stats from otto
type OttoStatsProject struct {
	ProjectID string
	Duration  time.Duration
}

func prjStats(w http.ResponseWriter, r *http.Request) {
	serial := r.FormValue("serial")
	otp := r.FormValue("otp")
	otto, err := db.Authorize(serial, otp)
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

	projectID := r.FormValue("project_id")
	if projectID == "" {
		log.Error("No valid project_id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	d := r.FormValue("date")
	log.Info(d)
	date, err := time.Parse("2006-01-02", d)
	if err != nil {
		log.Error("Wrong date format.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stats := OttoStatsProject{}
	entry := []db.ProjectEntry{}

	switch r.FormValue("when") {
	case "day":
		log.Info("Day:", date)
		if db.DB.Where("user_id = ? AND project_id = ? AND cast(start_time as date) = ?", otto.UserID, projectID, d).Find(&entry).RecordNotFound() {
			log.Error("Invalid Project ID: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, item := range entry {
			stats.Duration += item.EndTime.Sub(item.StartTime)
		}
	case "week":
		log.Info("week:", date)
	case "month":
		log.Info("month:", date)
		_, m, _ := date.Date()
		if db.DB.Where("user_id = ? AND project_id = ? AND EXTRACT(MONTH FROM start_time) = ?", otto.UserID, projectID, m).Find(&entry).RecordNotFound() {
			log.Error("Invalid Project ID: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, item := range entry {
			stats.Duration += item.EndTime.Sub(item.StartTime)
		}

	default:
		log.Error("Invalid when parameter: ", r.FormValue("when"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stats.ProjectID = projectID
	stats.Duration = stats.Duration / 1000000000
	log.Info("Stas: projectid: ", stats.ProjectID, ", duration: ", stats.Duration)
	js, err := json.Marshal(stats)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//ListenAndServe ... Otto endpoints
func ListenAndServe(addr string) error {
	http.HandleFunc("/ping", pong)
	http.HandleFunc("/projects", listProjects)
	http.HandleFunc("/register", registerEntry)
	http.HandleFunc("/stats", prjStats)
	http.HandleFunc("/slack_api", botto.InteractiveEventHandler)

	if conf.UseHTTPS() {
		return http.ListenAndServeTLS(addr,
			conf.GetCertFilePath(),
			conf.GetKeyfilePath(), nil)
	}

	return http.ListenAndServe(addr, nil)
}
