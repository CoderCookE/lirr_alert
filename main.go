package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func main() {
	go startPolling()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/lirr", Lirr)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Lirr(w http.ResponseWriter, r *http.Request) {
	fmt.Sprintf("http://www.movable-ink-7158.com/p/rp/60e6a4c03b713777.png?cache_buster=%s", time.Now().String())
}

type AlertChecker struct {
	Raw   string
	lines map[string]*Line
}

type Line struct {
	name   string
	status string
	text   string
}

func startPolling() {
	a := AlertChecker{}
	a.lines = make(map[string]*Line)

	for {
		<-time.After(15 * time.Second)
		go a.CheckAlert()
	}
}

func (a *AlertChecker) CheckAlert() {
	resp, err := http.Get("http://web.mta.info/status/serviceStatus.txt")
	if err != nil {
		log.Fatal(err.Error)
	}
	defer resp.Body.Close()
	rawData, _ := ioutil.ReadAll(resp.Body)
	htmlData := string(rawData)
	nameRegExp, _ := regexp.Compile("<name>(.*)</name>")
	statusRegExp, _ := regexp.Compile("<status>(.*)</status>")
	textRegExp, _ := regexp.Compile("<text>(.*)</text>")
	if htmlData != a.Raw {
		a.Raw = htmlData
		rawTrains := fmt.Sprintf("%q", strings.SplitAfter(htmlData, "LIRR>")[1])
		s := strings.SplitAfter(rawTrains, "</line>")

		for _, v := range s {
			var currentText string
			var currentStatus string

			if a.lines[v] != nil {
				currentText = a.lines[v].text
				currentStatus = a.lines[v].status
			}

			status := string(statusRegExp.Find([]byte(v)))
			text := string(textRegExp.Find([]byte(v)))
			name := string(nameRegExp.Find([]byte(v)))

			if currentStatus != status || currentText != text {

				log.Printf("%s", a.lines[v])

				a.lines[v] = &Line{name, status, text}

				log.Printf(a.lines[v].name, a.lines[v].text, a.lines[v].status)
			}
		}
	}
}
