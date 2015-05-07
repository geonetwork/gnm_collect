package gnmsys

import (
	"time"
	"net/http/cookiejar"
	"net/http"
	"log"
	"net/url"
	"strings"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type System interface {
	// function that will signal to the system to clean up and shutdown
	signalTerm()
	// function that will write the reports to disk
	signalFlush()
}

type SystemSignal int

const (
	term = 1 + iota
	flush
	tick
)

type defaultSystem struct {
	signals chan SystemSignal
	client  http.Client
	reports []Report
}

func CreateSystem(reports ...Report) defaultSystem {
	options := cookiejar.Options{}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	system := defaultSystem{
		reports: reports,
		signals: make(chan SystemSignal),
		client: http.Client{Jar: jar}}

	go loop(system.signals)

	return system
}

func loop(signals chan <- SystemSignal) {
	for {
		signals <- tick
		time.Sleep(time.Second)
	}
}

func (sys defaultSystem) signalTerm() {
	sys.signals <- term
}
func (sys defaultSystem) signalFlush() {
	sys.signals <- flush
}
func (sys defaultSystem) Run() {

	host := "http://tc-geocat.dev.bgdi.ch/geonetwork"

	log.Printf("Start Login \n")
	values := url.Values{"username":[]string{"testjesse"}, "password":[]string{"testjesse"}}

	resp, _ := sys.client.PostForm(host+"/j_spring_security_check", values)

	log.Printf("Login response: %q '%v': \n\n", resp.Status, resp.StatusCode)
	if resp.StatusCode > 300 {
		loc, _ := resp.Location()
		if loc == nil && !strings.Contains(loc.Path, "home") {
			log.Printf("Error %v", loc.Path)
			log.Fatalf("Error logging in: %q: '%v'\n", resp.Status, resp.StatusCode)
		}
	}

	var x int64 = 0

	for sig := range sys.signals {
		switch sig {
		case term:
			goto shutdown
		case flush:
			fmt.Printf("Not yet implemented")
		case tick:

		}
		resp, _ = sys.client.Get(host+"/monitor/metrics")
		log.Printf("Metrics response: %q '%v'\n", resp.Status, resp.StatusCode)
		if resp.StatusCode > 300 {
			log.Fatalf("Error obtaining metrics in: %q: '%v'\n", resp.Status, resp.StatusCode)
		}

		data, _ := ioutil.ReadAll(resp.Body)
		var jsonData map[string]interface{}

		json.Unmarshal(data, &jsonData)


		metrics := Json{jsonData}

		for _, report := range sys.reports {
			report.Update(x, metrics)
		}

		x++
	}

	shutdown:
	for _, report := range sys.reports {
		report.Save()
	}
}
