package gnmsys

import (
	"github.com/gonum/plot/plotter"
	"log"
)
type CollectorFactory func(capacity int) Collector

type Collector interface {
	Name() string
	AddSample(time int64, metrics Json)
	GetXYs() plotter.XYs
}

type FloatCollector struct {
	name           string
	capacity        int
	jsonPath       []string
	xys            plotter.XYs
}

func NewFloatCollector(name string, jsonPath ...string) CollectorFactory {
	fac := func(capacity int) Collector {
		return &FloatCollector{
			name: name,
			capacity: capacity,
			jsonPath: jsonPath,
			xys: plotter.XYs{}}
	}

	return fac
}

func (c *FloatCollector) AddSample(time int64, metrics Json) {
	y := metrics.resolveFloat(c.jsonPath...)
	log.Printf("Adding (%v, %v) to %q\n", time, y, c.name)
	xy := make(plotter.XYs, 1)
	xy[0].X = float64(time)
	xy[0].Y = y

	xys := c.xys
	if (len(c.xys) == c.capacity) {
		log.Printf("Capacity reached. ")
		xys = xys[1:]
	}
	c.xys = append(xys, xy[0])

	log.Printf("Added (%v, %v) to %q.  Number of samples: %v\n", time, y, c.name, len(c.xys))
}

func (c FloatCollector) GetXYs() plotter.XYs {
	return c.xys
}
func (c FloatCollector) Name() string {
	return c.name
}
