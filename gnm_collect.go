package main

import (
	"gnm_collect/gnmsys"
	"os"
	"log"
	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
)

type Listener interface {
	start()
}

func main() {
	configureLogs()

	reports := loadReports()

	sys := gnmsys.CreateSystem(reports...)
	listener := gnmsys.CliListener{sys}

	go listener.Start()

	sys.Run()
}

func loadReports() []gnmsys.Report {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Metrics Trend"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	report := gnmsys.NewLineGraphReport("test", p, 12 * vg.Inch, 10 * vg.Inch,
		gnmsys.NewFloatCollector("Total Used Mem", "jvm", "memory", "totalUsed"),
		gnmsys.NewFloatCollector("File Descriptor Usage", "jvm", "fd_usage"))

	return []gnmsys.Report{report}
}

func configureLogs() {
	f, err := os.OpenFile("gnm_collect.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
}
