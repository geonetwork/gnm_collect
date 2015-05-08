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
	"os"
	"path"
	"path/filepath"
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
type SysConfig struct {
	UrlStem,
	Username,
	Password,
	OutputDir string
	SampleConfigs []SampleConfig
}
type defaultSystem struct {
	config SysConfig
	signals chan SystemSignal
	client  http.Client
	reports []Report
}

func CreateSystem(config SysConfig, reportFactories ...ReportFactory) defaultSystem {
	for _, conf := range config.SampleConfigs {
		conf.Validate()
	}
	if config.OutputDir != "" {
		os.MkdirAll(config.OutputDir, os.ModeDir)
	} else {
		config.OutputDir = "."
	}
	options := cookiejar.Options{}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	reports := make([]Report, len(reportFactories) * len(config.SampleConfigs))
	if (len(reports) == 0) {
		log.Fatalf("No reports are configured\nReport Factories: %d\nSampleConfigs: %v\n", len(reportFactories), config.SampleConfigs)
	}
	for i, fac := range reportFactories {
		for j, sConf := range config.SampleConfigs {
			if config.OutputDir != "" {
				sConf.DirName = path.Join(config.OutputDir, sConf.DirName)
			}
			reports[(i * len(config.SampleConfigs)) + j] = fac(sConf)
		}
	}
	system := defaultSystem{
		config: config,
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
	fmt.Printf("Finalizing reports and shutting down...\n")
	sys.signals <- term
}
func (sys defaultSystem) signalFlush() {
	fmt.Printf("Saving/Flushing Reports to disk\n")
	sys.signals <- flush
}
func (sys defaultSystem) Run() {

	urlStem := sys.config.UrlStem

	log.Printf("Start Login \n")
	values := url.Values{"username":[]string{sys.config.Username}, "password":[]string{sys.config.Password}}

	resp, _ := sys.client.PostForm(urlStem+"/j_spring_security_check", values)

	log.Printf("Login response: %q '%v': \n\n", resp.Status, resp.StatusCode)
	if resp.StatusCode > 300 {
		loc, _ := resp.Location()
		if loc == nil || !strings.Contains(loc.Path, "home") {
			log.Printf("Error %v", loc.Path)
			log.Fatalf("Error logging in: %q: '%v'\n", resp.Status, resp.StatusCode)
		}
	}

	var timeSeconds int64 = 0

	for sig := range sys.signals {
		switch sig {
		case term:
			goto shutdown
		case flush:
			sys.save()
		case tick:
			resp, _ = sys.client.Get(sys.config.UrlStem+"/monitor/metrics")
			log.Printf("Metrics response: %q '%v'\n", resp.Status, resp.StatusCode)
			if resp.StatusCode > 300 {
				log.Fatalf("Error obtaining metrics in: %q: '%v'\n", resp.Status, resp.StatusCode)
			}

			data, _ := ioutil.ReadAll(resp.Body)
			var jsonData map[string]interface{}

			err := json.Unmarshal(data, &jsonData)

			if err != nil {
				msg := "Metrics response was not valid json, this is likely because the login username/password are incorrect. %v\n\n"
				fmt.Printf(msg, "")
				log.Fatalf(msg, err.Error())
			}
			metrics := Json{jsonData}


			for _, report := range sys.reports {
				if timeToUpdate(timeSeconds, report) {
					report.Update(timeSeconds, metrics)
				}
			}

			timeSeconds++
		}
	}

	shutdown:
	sys.save()

	fmt.Printf("\nSystem has Cleanly shutdown\n\n[DONE]\n")
}

func timeToUpdate(timeSeconds int64, report Report) bool {
	interval := int64(report.GetUpdateInterval())
	timeNano := timeSeconds * int64(time.Second)
	return timeNano % interval == 0
}

func (sys defaultSystem) save() {
	for _, report := range sys.reports {
		report.Save()
	}
	outputDir, _ := filepath.Abs(sys.config.OutputDir)
	fmt.Printf("Reports have been written to disk: '%s'\n", outputDir)
}