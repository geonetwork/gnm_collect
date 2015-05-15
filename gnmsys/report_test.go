package gnmsys

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"io/ioutil"
	"github.com/gonum/plot/plotter"
	"bytes"
)

var _ = Describe("Report", func() {
	log.SetOutput(ioutil.Discard)
	Describe("Save and Load", func() {
		It("A report can be saved and restored", func() {
			report := LineGraphReport{collectors: []Collector{&FloatCollector{
				name: "coll1",
				capacity: 100,
				xys: plotter.XYs{
					{X: 11, Y:11}, {X: 12, Y:12}, {X: 13, Y:13},
					{X: 14, Y:14}, {X: 15, Y:15}, {X: 16, Y:16}}},
				&FloatCollector{
					name: "coll2",
					capacity: 100,
					xys: plotter.XYs{
						{X: 21, Y:21}, {X: 22, Y:22}, {X: 23, Y:23},
						{X: 24, Y:24}, {X: 25, Y:25}}}}}
			buffer := &bytes.Buffer{}
			report.saveState(buffer)

			loadedReport := LineGraphReport{collectors:[]Collector{
				&FloatCollector{capacity: 100, name: "coll1"},
				&FloatCollector{capacity: 100, name: "coll2"}}}

			loadedReport.doLoad("test", buffer)

			Expect(len(loadedReport.collectors)).To(Equal(2))
			Expect(len(loadedReport.collectors[0].GetXYs())).To(Equal(6))
			Expect(len(loadedReport.collectors[1].GetXYs())).To(Equal(5))
		})
	})
})
