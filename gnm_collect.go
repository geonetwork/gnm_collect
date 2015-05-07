package main

import (
	"gnm_collect/gnmsys"
	"os"
	"log"
	"github.com/gonum/plot/vg"
)

type Listener interface {
	start()
}

func main() {
	configureLogs()

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
		Password: "testjesse"}
}

func loadReports() []gnmsys.ReportFactory {
	return [] gnmsys.ReportFactory{
		gnmsys.LineReportFactoryBuilder{
			Title: "Memory",
			Filename: "memory",
			YAxis: "Bytes",
			X: 12 * vg.Inch,
			Y: 5 * vg.Inch,
			Collectors: []gnmsys.CollectorFactory{
				gnmsys.NewFloatCollector("Max", "jvm", "memory", "totalMax"),
				gnmsys.NewFloatCollector("Total Used", "jvm", "memory", "totalUsed"),
				gnmsys.NewFloatCollector("Heap Used", "jvm", "memory", "heapUsed")}}.ToRequestFactory(),
		gnmsys.LineReportFactoryBuilder{
			Title: "Resource Usage (%)",
			Filename: "resource_usage",
			YAxis: "%",
			X: 12 * vg.Inch,
			Y: 5 * vg.Inch,
			Collectors: []gnmsys.CollectorFactory{
				gnmsys.NewFloatCollector("Non Heap Mem Used", "jvm", "memory", "non_heap_usage"),
				gnmsys.NewFloatCollector("Heap Mem Used", "jvm", "memory", "heap_usage"),
				gnmsys.NewFloatCollector("File Descriptor Usage", "jvm", "fd_usage")}}.ToRequestFactory()}
}

func configureLogs() {
	f, err := os.OpenFile("gnm_collect.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
}
