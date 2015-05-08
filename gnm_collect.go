package main

import (
	"github.com/geonetwork/gnm_collect/gnmsys"
	"os"
	"log"
	"github.com/gonum/plot/vg"
	"time"
	"flag"
	"io"
	"io/ioutil"
	"fmt"
	"net/http"
	"github.com/geonetwork/gnm_collect/gnmserver"
)

type Listener interface {
	start()
}

var targetServer, username, password, out string
var port = flag.Int("port", 10100, "The port to start the http server on")
var logging = flag.Bool("logging", false, "Set to true to enable logging")
func main() {
	flag.StringVar(&targetServer, "target", "http://localhost:8080", "The url of the server to test.  Default [http://localhost:8080/geonetwork")
	flag.StringVar(&username, "user", "", "The user name of a user that has Monitor privileges on the target server")
	flag.StringVar(&password, "pass", "", "The user's password")
	flag.StringVar(&out, "out", "./gnm_reports", "The directory to write the reports to.")
	flag.Parse()
	_, err := http.Get(fmt.Sprintf("http://localhost:%d/status", *port))

	if err == nil {
		fmt.Printf("Monitoring server is already running\n")
		os.Exit(0)
	}

	var logOut io.Writer = ioutil.Discard
	if *logging {
		f, err := os.OpenFile("gnm_collect.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		logOut = f
	}
	log.SetOutput(logOut)

	if targetServer[len(targetServer) - 1] == '/' {
		targetServer = targetServer[:len(targetServer) - 1]
	}

	config := loadSystemConfig()
	reports := loadReports()

	sys := gnmsys.CreateSystem(config, reports...)
	listener := gnmsys.CliListener{sys}

	go listener.Start()

	server := gnmserver.Server{
		Port: *port,
		Config: config,
		Sys: &sys}
	go server.Start()
	sys.Run()
}


func loadSystemConfig() gnmsys.SysConfig {
	if username == "" || password == "" {
		fmt.Fprintf(os.Stderr, "-user and -pass are required\n")
		flag.Usage()
		os.Exit(1)
	}
	return gnmsys.SysConfig{
		UrlStem: targetServer,
		Username: username,
		Password: password,
		OutputDir: out,
		SampleConfigs: []gnmsys.SampleConfig{
			gnmsys.SampleConfig{Name: "Last Five Minutes", DirName: "five_minutes", MaxSamples: 5 * 60, UpdateInterval: time.Second},
			gnmsys.SampleConfig{Name: "Last Five Hours", DirName: "five_hours", MaxSamples: 5 * 60, UpdateInterval: time.Minute},
			gnmsys.SampleConfig{Name: "Last Five Days", DirName: "five_days", MaxSamples: 5 * 24, UpdateInterval: time.Hour},
			gnmsys.SampleConfig{Name: "Last Five Months", DirName: "five_months", MaxSamples: 5 * 30, UpdateInterval: 24 * time.Hour},
		}}
}

func loadReports() []gnmsys.ReportFactory {
	return [] gnmsys.ReportFactory{
		gnmsys.LineReportFactoryBuilder{
			Title: "Memory",
			Filename: "memory",
			YAxis: "Bytes",
			X: 12 * vg.Inch,
			Y: 5 * vg.Inch,
			CollectorFactories: []gnmsys.CollectorFactory{
				gnmsys.NewFloatCollector("Max", "jvm", "memory", "totalMax"),
				gnmsys.NewFloatCollector("Total Used", "jvm", "memory", "totalUsed"),
				gnmsys.NewFloatCollector("Heap Used", "jvm", "memory", "heapUsed")}}.ToRequestFactory(),
		gnmsys.LineReportFactoryBuilder{
			Title: "Resource Usage (%)",
			Filename: "resource_usage",
			YAxis: "%",
			X: 12 * vg.Inch,
			Y: 5 * vg.Inch,
			CollectorFactories: []gnmsys.CollectorFactory{
				gnmsys.NewFloatCollector("Non Heap Mem Used", "jvm", "memory", "non_heap_usage"),
				gnmsys.NewFloatCollector("Heap Mem Used", "jvm", "memory", "heap_usage"),
				gnmsys.NewFloatCollector("File Descriptor Usage", "jvm", "fd_usage")}}.ToRequestFactory(),
		gnmsys.LineReportFactoryBuilder{
			Title: "Time Used by Garbage Collectors",
			Filename: "garbage_collectors",
			YAxis: "Milliseconds",
			X: 12 * vg.Inch,
			Y: 5 * vg.Inch,
			CollectorFactories: []gnmsys.CollectorFactory{
				gnmsys.NewFloatCollector("MarkSweep", "jvm", "garbage-collectors", "PS MarkSweep", "time"),
				gnmsys.NewFloatCollector("Scavenge", "jvm", "garbage-collectors", "PS Scavenge", "time")}}.ToRequestFactory(),
		gnmsys.LineReportFactoryBuilder{
			Title: "CPU Load",
			Filename: "cpu_load",
			YAxis: "Load (%)",
			X: 12 * vg.Inch,
			Y: 5 * vg.Inch,
			CollectorFactories: []gnmsys.CollectorFactory{
				gnmsys.NewFloatCollector("System", "java.lang.management.OperatingSystemMXBean",
					"Process_CPU_Load", "value"),
				gnmsys.NewFloatCollector("System", "java.lang.management.OperatingSystemMXBean",
					"Process_CPU_Load", "value"),
				gnmsys.NewFloatCollector("Average Load", "java.lang.management.OperatingSystemMXBean",
					"systemLoadAverage", "value")}}.ToRequestFactory(),
		gnmsys.LineReportFactoryBuilder{
			Title: "Thread states (%)",
			Filename: "thread_states",
			YAxis: "%",
			X: 12 * vg.Inch,
			Y: 5 * vg.Inch,
			CollectorFactories: []gnmsys.CollectorFactory{
				gnmsys.NewFloatCollector("Blocked", "jvm", "thread-states", "blocked"),
				gnmsys.NewFloatCollector("New", "jvm", "thread-states", "new"),
				gnmsys.NewFloatCollector("Timed Waiting", "jvm", "thread-states", "timed_waiting"),
				gnmsys.NewFloatCollector("Waiting", "jvm", "thread-states", "waiting"),
				gnmsys.NewFloatCollector("Terminated", "jvm", "thread-states", "terminated"),
				gnmsys.NewFloatCollector("Runnable", "jvm", "thread-states", "runnable")}}.ToRequestFactory()}
}

