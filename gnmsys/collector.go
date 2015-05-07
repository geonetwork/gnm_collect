package gnmsys

import (
	"github.com/gonum/plot/plotter"
	"log"
)

type Collector interface {
	Name() string
	AddSample(time int64, metrics Json)
	GetXYs() plotter.XYs
}

type FloatCollector struct {
	name           string
	jsonPath       []string
	xys            plotter.XYs
}

func NewFloatCollector(name string, jsonPath ...string) Collector {
	return &FloatCollector{
		name: name,
		jsonPath: jsonPath,
		xys: plotter.XYs{}}
}

func (c *FloatCollector) AddSample(time int64, metrics Json) {
	y := metrics.resolveFloat(c.jsonPath...)
	log.Printf("Adding (%v, %v) to %q\n", time, y, c.name)
	xy := make(plotter.XYs, 1)
	xy[0].X = float64(time)
	xy[0].Y = y

	c.xys = append(c.xys, xy[0])

	log.Printf("Added (%v, %v) to %q.  Number of samples: %v\n", time, y, c.name, len(c.xys))
}

func (c FloatCollector) GetXYs() plotter.XYs {
	return c.xys
}
func (c FloatCollector) Name() string {
	return c.name
}
