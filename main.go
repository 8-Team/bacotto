package main

import (
	"flag"
	"time"

	"github.com/8-team/bacotto/api"
	"github.com/8-team/bacotto/botto"
	"github.com/8-team/bacotto/conf"
	"github.com/8-team/bacotto/db"
	"github.com/8-team/bacotto/erp"
	"github.com/8-team/bacotto/slackbot"
	log "github.com/Sirupsen/logrus"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var cfgpath = flag.String("conf", "res/default.toml", "configuration file path")

func main() {
	flag.Parse()

	if err := conf.Load(*cfgpath); err != nil {
		log.Warnln("Configuration file not found, falling back to defaults")
	}

	if err := erp.Open(); err != nil {
		panic(err)
	}
	defer erp.Close()

	if err := db.Open(conf.GetDatabaseURI()); err != nil {
		panic(err)
	}
	defer db.Close()

	if conf.GetSerialsURI() != "" {
		if err := db.SyncSerial(conf.GetSerialsURI()); err != nil {
			panic(err)
		}
	}

	bot := slackbot.New(conf.GetSlackToken())

	bot.OnContextCreation(botto.CheckUserPresence)

	bot.Register("add project", "start managing a project using your Otto", botto.PickProject)
	bot.Register("list projects", "show a list of your managed projects", botto.ListProjects)
	bot.Register("report", "show a report of today's activities", botto.ShowReport)

	go func() {

		for {
			if err := bot.Start(); err != nil {
				log.WithField("app", "botto").Errorln(err)
			}
			log.Warningln("Starting again in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}()

	if err := api.ListenAndServe(conf.GetHTTPListenAddr(), bot.InteractiveEventHandler); err != nil {
		log.WithField("app", "api").Fatalln(err)
	}
}
