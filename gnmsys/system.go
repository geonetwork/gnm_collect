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
	"io"
	"github.com/gonum/plot"
	"runtime/debug"
)

type System interface {
	// function that will signal to the system to clean up and shutdown
	SignalTerm()
	// function that will write the reports to disk
	SignalFlush()
	// List all the Reports available
	GetReports() []Report
	// Get the output file of the report
	GetReportFile(report Report) string
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
	if config.OutputDir == "" || config.OutputDir == "." {
		config.OutputDir = "gnm_reports"
	} else {
		os.MkdirAll(config.OutputDir, os.FileMode(0755))
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

	system.validate()

	go loop(system.signals)

	return system
}

func (sys defaultSystem) validate() {
	for _, conf := range sys.config.SampleConfigs {
		conf.Validate()
	}

	probeFile := filepath.Join(sys.config.OutputDir, "probe")
	if err := ioutil.WriteFile(probeFile, []byte("t"), os.FileMode(0664)); err != nil {
		log.Fatalf("Do not have write permissions to %s\n", sys.config.OutputDir)
	}
	os.Remove(probeFile)

	if _, err := plot.New(); err != nil {
		fmt.Printf("Error creating a test graph %q\n", err.Error())
		msg := "A likely problem is that the directory with font files cannot be found." +
		"  Copy https://github.com/gonum/plot/tree/master/vg/fonts to same directory as the executable"
		fmt.Printf(msg)

		log.Printf("Error creating a test graph %q\n", err.Error())
		log.Fatalf(msg)
	}
}

func loop(signals chan <- SystemSignal) {
	for {
		signals <- tick
		time.Sleep(time.Second)
	}
}

func (sys defaultSystem) GetReportFile(report Report) string {
	var catDirName string
	for _, sampConf := range sys.config.SampleConfigs {
		if sampConf.Name == report.GetCategory() {
			catDirName = sampConf.DirName
			break;
		}
	}
	return filepath.Join(sys.config.OutputDir, catDirName, report.GetFileName())
}
func (sys defaultSystem) GetReports() []Report {
	return sys.reports
}
func (sys defaultSystem) SignalTerm() {
	fmt.Printf("Finalizing reports and shutting down...\n")
	sys.signals <- term
}
func (sys defaultSystem) SignalFlush() {
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
		fmt.Println("System has started")
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
			msg := "Recovering from Panic: %v\n"
			fmt.Printf(msg, r)
			log.Printf(msg, r)
			debug.PrintStack()
			state.mustLogin = false
		}
	}()

	if (state.mustLogin) {
		loginUrl := state.urlStem+"/j_spring_security_check"
		log.Printf("Start Login: %s \n", loginUrl)
		resp, _ := sys.client.PostForm(loginUrl, state.loginCredentials)

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
	metricsUrl := sys.config.UrlStem+"/monitor/metrics"
	log.Printf("Making Metrics request %s", metricsUrl)
	resp, _ := sys.client.Get(metricsUrl)
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

	for i, report := range sys.reports {
		log.Printf("Check if report should be updated: %q\n", report.GetName())
		log.Printf("report %d of %d\n", i, len(sys.reports))
		if timeToUpdate(int64(requestTime), report) {
			safeUpdateReport(report, metrics, requestTime)
		}
	}

	if timeToWriteGraphs(requestTime, state.startTime) {
		sys.save(state.startTime.Format(timeFmt))
	}
}

func safeUpdateReport(report Report, metrics Json, requestTime int64) {
	defer func() {
		if r := recover(); r != nil {
			msg := "Recovering from error in safeUpdateReport\n    Report %q\n    Metrics: %v\n    Error %v\n"
			fmt.Printf(msg, report.GetName(), metrics.Data, r)
			log.Printf(msg, r)
		}
	}()
	report.Update(requestTime, metrics)
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
	tmpDir := path.Join(sys.config.OutputDir, "tmp")
	for _, report := range sys.reports {
		report.Save(titleModifier, tmpDir)
	}

	filepath.Walk(tmpDir, func(file string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			rel, err := filepath.Rel(tmpDir, file)
			if err == nil {
				dest := path.Join(sys.config.OutputDir, rel)
				os.Remove(dest)
				if _, err := os.Stat(dest); os.IsNotExist(err) {
					os.MkdirAll(filepath.Dir(dest), os.FileMode(0755))
					err = os.Rename(file, dest)
					if err == nil {
						log.Printf("Moved %s to %s\n", file, dest)
					} else {
						log.Printf("Error occurred when attempting to move %s to %s:\n%q\n", file, dest, err.Error())
						copyMove(file, dest)
					}
				} else {
					copyMove(file, dest)
				}
			}
		}
		return nil
	})

	os.RemoveAll(tmpDir)
}

func copyMove(source, dest string){
	sFile, err := os.Create(source)

	if err != nil {
		log.Printf("Failed to open/create source file %s in copy@sytem.go:\n%q\n", source, err.Error())
		return
	}
	defer sFile.Close()

	dFile, err := os.Create(dest)
	if err != nil {
		log.Printf("Failed to open/create dest file %s in copy@sytem.go:\n%q\n", dest, err.Error())
		return
	}
	defer dFile.Close()

	if _, err = io.Copy(dFile, sFile); err != nil {
		log.Printf("Error occurred trying to copy %s to %s:\n%q\n", source, dest, err.Error())
		return
	}
	os.Remove(source)
}