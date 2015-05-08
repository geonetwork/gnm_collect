package gnmsys
import (
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot"
)

type ReportFactory func(maxSamples int) Report

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

	return func(maxSamples int) Report {
		collectors := make([]Collector, len(b.CollectorFactories))
		for i, cFactory := range b.CollectorFactories {
			collectors[i] = cFactory(maxSamples)
		}
		return NewLineGraphReport(b.Filename, newPlot, b.X, b.Y, collectors...)
	}
}