package main

import (
	"gnm_collect/gnmsys"
	"os"
	"log"
	"github.com/gonum/plot/vg"
	"time"
)

type Listener interface {
	start()
}

func main() {
	f, err := os.OpenFile("gnm_collect.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	config := loadSystemConfig()
	reports := loadReports()

	sys := gnmsys.CreateSystem(config, reports...)
	listener := gnmsys.CliListener{sys}

	go listener.Start()

	sys.Run()
}

func loadSystemConfig() gnmsys.SysConfig {
	return gnmsys.SysConfig{
		UrlStem: "http://tc-geocat.dev.bgdi.ch/geonetwork",
		Username: "testjesse",
		Password: "testjesse",
		OutputDir: "./reports",
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

