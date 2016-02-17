package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type AlertChecker struct {
	Raw   string
	lines map[string]*Line
}

func (a *AlertChecker) CheckAlert() {
	resp, err := http.Get("http://web.mta.info/status/serviceStatus.txt")
	if err != nil {
		log.Fatal(err.Error)
	}
	defer resp.Body.Close()
	rawData, _ := ioutil.ReadAll(resp.Body)
	htmlData := string(rawData)
	if htmlData != a.Raw {
		a.Raw = htmlData
		rawTrains := fmt.Sprintf("%q", strings.SplitAfter(htmlData, "LIRR>")[1])
		s := strings.SplitAfter(rawTrains, "</line>")

		for _, v := range s {
			go a.CheckLine(v)
		}
	}
}

func (a *AlertChecker) CheckLine(v string) {
	nameRegExp, _ := regexp.Compile("<name>(.*)</name>")
	statusRegExp, _ := regexp.Compile("<status>(.*)</status>")
	textRegExp, _ := regexp.Compile("<text>(.*)</text>")
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
		a.lines[v] = &Line{name, status, text}
		pushToSlack(a.lines[v].ToString())
	}
}

func pushToSlack(message string) {
	log.Printf(message)
}
