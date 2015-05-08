package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gnm_collect/gnmsys"
	"log"
	"io/ioutil"
)

var _ = Describe("FloatCollector", func() {
	log.SetOutput(ioutil.Discard)
		Describe("Add Points to Float Collector", func() {
			Context("Empty collection", func() {
					It("Adding 3 samples should result in a xyz with 3 samples", func () {
							coll := gnmsys.NewFloatCollector("Test", "data")(200)
							coll.AddSample(1, gnmsys.Json{map[string]interface {}{"data" : 32}})
							coll.AddSample(2, gnmsys.Json{map[string]interface {}{"data" : 34}})
							coll.AddSample(3, gnmsys.Json{map[string]interface {}{"data" : 36}})

							Expect (3).To(Equal(len(coll.GetXYs())))

							Expect (1.0).To(Equal(coll.GetXYs()[0].X))
							Expect (32.0).To(Equal(coll.GetXYs()[0].Y))

							Expect (2.0).To(Equal(coll.GetXYs()[1].X))
							Expect (34.0).To(Equal(coll.GetXYs()[1].Y))

							Expect (3.0).To(Equal(coll.GetXYs()[2].X))
							Expect (36.0).To(Equal(coll.GetXYs()[2].Y))

						})
				})
				Context("Full Collector", func() {
					It ("Adding 1 sample to full collector will drop oldest value and add last sample as last value", func() {
						coll := gnmsys.NewFloatCollector("test", "data")(3)

						coll.AddSample(1, gnmsys.Json{map[string]interface {}{"data" : 32}})
						coll.AddSample(2, gnmsys.Json{map[string]interface {}{"data" : 34}})
						coll.AddSample(3, gnmsys.Json{map[string]interface {}{"data" : 36}})

						Expect (len(coll.GetXYs())).To(Equal(3))

						Expect (1.0).To(Equal(coll.GetXYs()[0].X))
						Expect (3.0).To(Equal(coll.GetXYs()[2].X))

						coll.AddSample(4, gnmsys.Json{map[string]interface {}{"data" : 40}})

						Expect (3).To(Equal(len(coll.GetXYs())))

						Expect (2.0).To(Equal(coll.GetXYs()[0].X))
						Expect (4.0).To(Equal(coll.GetXYs()[2].X))
						Expect (40.0).To(Equal(coll.GetXYs()[2].Y))
					})
				})
			})
	})
