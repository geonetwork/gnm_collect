package gnmsys
import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
)

type Report interface {
	Update(time int64, metrics Json)
	Save()
}

type LineGraphReport struct {
	collectors []Collector
	name string
	report *plot.Plot
	width, height vg.Length
}

func NewLineGraphReport(name string,
						report *plot.Plot,
						width, height vg.Length,
						collectors ...Collector) LineGraphReport {
	return LineGraphReport{collectors: collectors, name: name, report: report, width:width, height:height}
}

func (r LineGraphReport) Update(time int64, metrics Json) {
	for _, coll := range r.collectors {
		coll.AddSample(time, metrics)
	}
}

func (r LineGraphReport) Save() {
	for _, coll := range r.collectors {
		err := plotutil.AddLinePoints(r.report, coll.Name(), coll.GetXYs())
		if err != nil {
			panic(err)
		}
	}

	if err := r.report.Save(r.width, r.height, r.name + ".png"); err != nil {
		panic(err)
	}
}