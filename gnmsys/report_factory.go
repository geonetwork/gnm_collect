package gnmsys
import (
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot"
	"time"
	"log"
	"github.com/geonetwork/gnm_collect/gnmsys/unit"
	"io/ioutil"
	"path/filepath"
	"os"
)

type SampleConfig struct {
	Name, DirName, XAxis string
	MaxSamples int
	UpdateInterval time.Duration
}
func(conf SampleConfig) Unit() unit.Unit {
	return unit.FindUnit(conf.UpdateInterval)
}
func(conf SampleConfig) Validate() {
	if conf.MaxSamples < 1 {
		log.Fatalf("Sample Config %q is invalid. MaxSamples must be > 0: %d", conf.Name, conf.MaxSamples)
	}
	if conf.UpdateInterval < time.Second {
		log.Fatalf("Sample Config %q is invalid. UpdateInteral must be >= 1 Second: %d", conf.Name, conf.UpdateInterval)
	}
}


type ReportFactory func(sampleConfig SampleConfig) Report

type LineReportFactoryBuilder struct {
	Title, YAxis, Filename string
	X, Y vg.Length
	CollectorFactories []CollectorFactory
}
func (b LineReportFactoryBuilder) ToRequestFactory() ReportFactory {
	newPlot := func() *plot.Plot {
		p, err := plot.New()
		if err != nil {
			panic(err)
		}

		p.Title.Text = b.Title
		p.Y.Label.Text = b.YAxis
		return p
	}

	return func(sampleConfig SampleConfig) Report {
		collectors := make([]Collector, len(b.CollectorFactories))
		for i, cFactory := range b.CollectorFactories {
			collectors[i] = cFactory(sampleConfig.MaxSamples)
		}
		return NewLineGraphReport(b.Filename, newPlot, b.X, b.Y, sampleConfig, collectors...)
	}
}