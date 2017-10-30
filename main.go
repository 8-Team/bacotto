package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/8-team/bacotto/botto"
	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const port = 4273

func getDbURI() string {
	return os.Getenv("DB_URI")
}

func getSlackToken() string {
	return os.Getenv("BOTTO_API_TOKEN")
}

// The API of bacotto
type API struct {
	// unused for now
	db *gorm.DB
}

func (api *API) helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Println("profit!")
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	db, err := gorm.Open("postgres", getDbURI())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	go func() {
		for {
			if err := botto.ListenAndServe(getSlackToken()); err != nil {
				log.WithField("app", "botto").Errorln(err)
				log.Warningln("Retrying in 1 second...")
				time.Sleep(time.Second)
			}
		}
	}()

	api := API{db}

	http.HandleFunc("/", api.helloWorld)

	addr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(addr, nil)
}
