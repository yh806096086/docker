package server

import (
	"net/http"
	"time"
	"log"
)


func init() {

}

func Run() {
	s := http.Server{
		Addr: ":8080",
		Handler: nil,
		ReadTimeout: 1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
