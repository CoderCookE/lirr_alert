package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	go startPolling()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/lirr", Lirr)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Lirr(w http.ResponseWriter, r *http.Request) {
	cachebuster := rand.Intn(50)
	w.Write([]byte(fmt.Sprintf("http://www.movable-ink-7158.com/p/rp/60e6a4c03b713777.png?cache_buster=%v", cachebuster)))
}

func startPolling() {
	a := AlertChecker{}
	a.lines = make(map[string]*Line)
	a.CheckAlert()

	for {
		<-time.After(60 * time.Second)
		go a.CheckAlert()
	}
}
