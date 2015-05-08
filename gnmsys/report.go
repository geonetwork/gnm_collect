package gnmsys
import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"path"
	"os"
	"log"
	"time"
)

type PlotFactory func() *plot.Plot

type Report interface {
	GetUpdateInterval() time.Duration
	Update(timeSeconds int64, metrics Json)
	Save()
}

type LineGraphReport struct {
	collectors []Collector
	name string
	report PlotFactory
	width, height vg.Length
	sampleConfig SampleConfig
}

func NewLineGraphReport(name string,
report PlotFactory,
width, height vg.Length,
sampleConfig SampleConfig,
collectors ...Collector) LineGraphReport {
	return LineGraphReport{collectors: collectors, name: name, report: report, width:width, height:height, sampleConfig: sampleConfig}
}

func (r LineGraphReport) Update(timeSeconds int64, metrics Json) {
	for _, coll := range r.collectors {
		coll.AddSample(r.sampleConfig.Unit().ConvertSeconds(timeSeconds), metrics)
	}
}

func (r LineGraphReport) GetUpdateInterval() time.Duration {
	return r.sampleConfig.UpdateInterval
}
func (r LineGraphReport) Save() {
	report := r.report()
	report.X.Label.Text = r.sampleConfig.Unit().String()
	lines := make([]interface{}, len(r.collectors) * 2)
	for i, coll := range r.collectors {
		lines[i * 2] = coll.Name()
		lines[i * 2 + 1] = coll.GetXYs()
	}

	err := plotutil.AddLinePoints(report, lines...)

	if err != nil {
		panic(err)
	}

	if (r.sampleConfig.DirName != "") {
		os.MkdirAll(r.sampleConfig.DirName, os.ModeDir)
	}
	outputFile := path.Join(r.sampleConfig.DirName, r.name) + ".png"
	log.Printf("Saving %q to %q", r.name, outputFile)
	if err := report.Save(r.width, r.height, outputFile); err != nil {
		panic(err)
	}
}