package api

import (
	"net/http"

	"github.com/8-team/bacotto/conf"
)

func pong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("HTTPS pong"))
}

func ListenAndServe(addr string) error {
	http.HandleFunc("/ping", pong)

	return http.ListenAndServeTLS(addr,
		conf.GetCertFilePath(),
		conf.GetKeyfilePath(), nil)
}
