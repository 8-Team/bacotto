package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const port = 4273

func getDbURI() string {
	return os.Getenv("DB_URI")
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

	api := API{db}

	http.HandleFunc("/", api.helloWorld)

	addr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(addr, nil)
}
