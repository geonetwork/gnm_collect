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
	GetCategory() string
	GetName() string
	GetFileName() string
	GetUpdateInterval() time.Duration
	Update(timeSeconds int64, metrics Json)
	Save(titleModifier string, reportDir string)
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
func (r LineGraphReport) GetCategory() string {
	return r.sampleConfig.Name
}
func (r LineGraphReport) GetName() string {
	return r.name
}
func (r LineGraphReport) GetFileName() string {
	return r.name + ".png"
}
func (r LineGraphReport) Update(timeSeconds int64, metrics Json) {
	log.Printf("Updating report: %q\n", r.GetName())

	for _, coll := range r.collectors {
		coll.AddSample(r.sampleConfig.Unit().ConvertSeconds(timeSeconds), metrics)
	}
}

func (r LineGraphReport) GetUpdateInterval() time.Duration {
	return r.sampleConfig.UpdateInterval
}
func (r LineGraphReport) Save(titleModifier string, reportDir string) {
	report := r.report()
	report.X.Label.Text = r.sampleConfig.Unit().String()
	report.Title.Text = report.Title.Text + "(" + titleModifier + ")"
	lines := make([]interface{}, len(r.collectors) * 2)
	for i, coll := range r.collectors {
		lines[i * 2] = coll.Name()
		lines[i * 2 + 1] = coll.GetXYs()
	}

	err := plotutil.AddLinePoints(report, lines...)

	if err != nil {
		panic(err)
	}
	outDir := path.Join(reportDir, r.sampleConfig.DirName)
	os.MkdirAll(outDir, os.FileMode(0755))

	outputFile := path.Join(outDir, r.GetFileName())
	log.Printf("Saving %q to %q", r.name, outputFile)
	if err := report.Save(r.width, r.height, outputFile); err != nil {
		panic(err)
	}
}