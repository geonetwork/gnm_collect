package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gnm_collect/gnmsys"
)

var _ = Describe("FloatCollector", func() {
		var coll = gnmsys.NewFloatCollector("Test", "data")

		Describe("Add Points to Float Collector", func() {
			Context("Empty collection", func() {
					It("Addint 3 samples should result in a xyz with 3 samples", func () {
							gnmsys.NewFloatCollector("Test", "data")
							coll.AddSample(1, gnmsys.Json{map[string]interface {}{"data" : 32}})
							coll.AddSample(2, gnmsys.Json{map[string]interface {}{"data" : 34}})
							coll.AddSample(3, gnmsys.Json{map[string]interface {}{"data" : 36}})

							Expect (len(coll.GetXYs())).To(Equal(3))

							Expect (coll.GetXYs()[0].X).To(Equal(float64(1)))
							Expect (coll.GetXYs()[0].Y).To(Equal(float64(32)))

							Expect (coll.GetXYs()[1].X).To(Equal(float64(2)))
							Expect (coll.GetXYs()[1].Y).To(Equal(float64(34)))

							Expect (coll.GetXYs()[2].X).To(Equal(float64(3)))
							Expect (coll.GetXYs()[2].Y).To(Equal(float64(36)))
						})
				})
			})
	})
