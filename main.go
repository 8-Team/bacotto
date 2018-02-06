package main

import (
	"flag"
	"time"

	"github.com/8-team/bacotto/api"
	"github.com/8-team/bacotto/botto"
	"github.com/8-team/bacotto/conf"
	"github.com/8-team/bacotto/db"
	"github.com/8-team/bacotto/erp"
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

	go func() {
		for {
			if err := botto.ListenAndServe(conf.GetSlackToken()); err != nil {
				log.WithField("app", "botto").Errorln(err)
			}
			log.Warningln("Starting again in 1 second...")
			time.Sleep(time.Second)
		}
	}()

	if err := api.ListenAndServe(conf.GetHTTPListenAddr()); err != nil {
		log.WithField("app", "api").Fatalln(err)
	}
}
