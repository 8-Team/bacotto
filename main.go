package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/8-team/bacotto/botto"
	"github.com/8-team/bacotto/db"
	log "github.com/Sirupsen/logrus"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const port = 4273

func getDbURI() string {
	return os.Getenv("DB_URI")
}

func getSlackToken() string {
	return os.Getenv("BOTTO_API_TOKEN")
}

func main() {
	db, err := db.Open(getDbURI())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	go func() {
		for {
			if err := botto.ListenAndServe(getSlackToken()); err != nil {
				log.WithField("app", "botto").Errorln(err)
			}
			log.Warningln("Starting again in 1 second...")
			time.Sleep(time.Second)
		}
	}()

	addr := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.WithField("app", "api").Fatalln(err)
	}
}
