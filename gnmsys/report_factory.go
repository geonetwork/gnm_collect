package gnmsys
import (
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot"
)

type ReportFactory func() Report

type LineReportFactoryBuilder struct {
	Title, YAxis, Filename string
	X, Y vg.Length
	Collectors []CollectorFactory
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

	return func() Report {
		collectors := make([]Collector, len(b.Collectors))
		for i, c := range b.Collectors {
			collectors[i] = c()
		}
		return NewLineGraphReport(b.Filename, newPlot, b.X, b.Y, collectors...)
	}
}