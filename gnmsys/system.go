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
	if config.OutputDir == "" || config.OutputDir == "." {
		config.OutputDir = "reports"
	} else {
		os.MkdirAll(config.OutputDir, os.ModeDir)
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

type systemState struct {
	startTime time.Time
	initializationComplete, mustLogin bool
	urlStem string
	loginCredentials url.Values
}
const timeFmt = "Start Time: 2006 Jan _2 15:04:05"
func (state *systemState) initialize() {
	if !state.initializationComplete {
		state.initializationComplete = true
		state.startTime = time.Now()
	}
}
func (sys defaultSystem) Run() {

	state := &systemState{
		urlStem: sys.config.UrlStem,
		loginCredentials: url.Values{"username":[]string{sys.config.Username}, "password":[]string{sys.config.Password}},
		initializationComplete: false,
		mustLogin: true}

	for sig := range sys.signals {
		switch sig {
		case term:
			goto shutdown
		case flush:
			sys.save(state.startTime.Format(timeFmt))
		case tick:
			sys.pollMetrics(state)
		}
	}

	shutdown:
	sys.save(state.startTime.Format(timeFmt))

	fmt.Printf("\nSystem has Cleanly shutdown\n\n[DONE]\n")
}

func (sys defaultSystem) pollMetrics(state *systemState) {
	defer func() {
		if r := recover(); r != nil {
			state.mustLogin = false
		}
	}()

	if (state.mustLogin) {
		log.Printf("Start Login \n")
		resp, _ := sys.client.PostForm(state.urlStem+"/j_spring_security_check", state.loginCredentials)

		log.Printf("Login response: %q '%v': \n\n", resp.Status, resp.StatusCode)
		state.mustLogin = false
		if resp.StatusCode > 300 {
			loc, _ := resp.Location()
			if loc == nil || !strings.Contains(loc.Path, "home") {
				log.Panicf("Error %v", loc.Path)
			}
		}
	}
	state.initialize()

	requestTime := time.Now().Unix() - state.startTime.Unix()
	resp, _ := sys.client.Get(sys.config.UrlStem+"/monitor/metrics")
	log.Printf("Metrics response: %q '%v'\n", resp.Status, resp.StatusCode)
	if resp.StatusCode > 300 {
		log.Panicf("Error obtaining metrics in: %q: '%v'\n", resp.Status, resp.StatusCode)
	}

	data, _ := ioutil.ReadAll(resp.Body)
	var jsonData map[string]interface{}

	err := json.Unmarshal(data, &jsonData)

	if err != nil {
		msg := "Metrics response was not valid json %v\n\n"
		log.Panicf(msg, err.Error())
	}
	metrics := Json{jsonData}

	for _, report := range sys.reports {
		if timeToUpdate(int64(requestTime), report) {
			report.Update(int64(requestTime), metrics)
		}
	}

	if timeToWriteGraphs(requestTime, state.startTime) {
		sys.save(state.startTime.Format(timeFmt))
	}
}

func timeToWriteGraphs(requestTime int64, startTime time.Time) bool {
	timeDiff := (time.Now().Second() - startTime.Second())
	return requestTime > 60 && timeDiff == 0
}
func timeToUpdate(timeSeconds int64, report Report) bool {
	interval := int64(report.GetUpdateInterval())
	timeNano := timeSeconds * int64(time.Second)
	return timeNano % interval == 0
}

func (sys defaultSystem) save(titleModifier string) {
	tmpDir := path.Join(os.TempDir(), "gnm_collect_tmp")
	for _, report := range sys.reports {
		report.Save(titleModifier, tmpDir)
	}
	os.RemoveAll(sys.config.OutputDir)
	os.Rename(tmpDir, sys.config.OutputDir)
	outputDir, _ := filepath.Abs(sys.config.OutputDir)
	fmt.Printf("Reports have been written to disk: '%s'\n", outputDir)
}