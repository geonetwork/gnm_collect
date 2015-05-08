package gnmsys

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"io/ioutil"
	"time"
)

var _ = Describe("System", func() {
	log.SetOutput(ioutil.Discard)
	Describe("timeToUpdate Method", func() {
		It("Update each time with the report duration is each second", func() {
			report := LineGraphReport{sampleConfig: SampleConfig{UpdateInterval: time.Second}}
			Expect(timeToUpdate(0, report)).To(BeTrue())
			Expect(timeToUpdate(1, report)).To(BeTrue())
		})
		It("Update each time with the report duration is each minute", func() {
			report := LineGraphReport{sampleConfig: SampleConfig{UpdateInterval: time.Minute}}
			Expect(timeToUpdate(0, report)).To(BeTrue())
			Expect(timeToUpdate(1, report)).To(BeFalse())
			Expect(timeToUpdate(3, report)).To(BeFalse())
			Expect(timeToUpdate(60, report)).To(BeTrue())
			Expect(timeToUpdate(64, report)).To(BeFalse())
			Expect(timeToUpdate(120, report)).To(BeTrue())
		})
		It("Update each time with the report duration is each hour", func() {
			report := LineGraphReport{sampleConfig: SampleConfig{UpdateInterval: time.Hour}}
			Expect(timeToUpdate(0, report)).To(BeTrue())
			Expect(timeToUpdate(1, report)).To(BeFalse())
			Expect(timeToUpdate(3, report)).To(BeFalse())
			Expect(timeToUpdate(60, report)).To(BeFalse())
			Expect(timeToUpdate(64, report)).To(BeFalse())
			Expect(timeToUpdate(60 * 60, report)).To(BeTrue())
			Expect(timeToUpdate(60 * 60 + 4, report)).To(BeFalse())
			Expect(timeToUpdate(2 * 60 * 60 + 59, report)).To(BeFalse())
		})
	})
})
