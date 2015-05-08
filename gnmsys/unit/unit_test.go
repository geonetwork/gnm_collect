package unit

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"io/ioutil"
	"testing"
	"time"
)


func TestPrivateGnmsysAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unit Test Suite")
}

var _ = Describe("Unit", func() {
	log.SetOutput(ioutil.Discard)
	Describe("convertSeconds Method", func() {
		It("should convert a unit in seconds to a new value in the new unit", func() {
			Expect(Seconds.ConvertSeconds(1)).To(Equal(int64(1)))
			Expect(Seconds.ConvertSeconds(100)).To(Equal(int64(100)))
			Expect(Days.ConvertSeconds(1)).To(Equal(int64(0)))
			Expect(Days.ConvertSeconds(24)).To(Equal(int64(0)))
			Expect(Days.ConvertSeconds(60 * 60 * 24)).To(Equal(int64(1)))
			Expect(Days.ConvertSeconds(2 * 60 * 60 * 24)).To(Equal(int64(2)))

			Expect(Minutes.ConvertSeconds(1)).To(Equal(int64(0)))
			Expect(Minutes.ConvertSeconds(24)).To(Equal(int64(0)))
			Expect(Minutes.ConvertSeconds(60)).To(Equal(int64(1)))
			Expect(Minutes.ConvertSeconds(62)).To(Equal(int64(1)))
			Expect(Minutes.ConvertSeconds(120)).To(Equal(int64(2)))
			Expect(Minutes.ConvertSeconds(122)).To(Equal(int64(2)))
		})
	})
	Describe("FindUnit() Method", func() {
		FIt("Should calculate the correct unit for the given updateInterval", func() {
			Expect(FindUnit(time.Second)).To(Equal(Seconds))
			Expect(FindUnit(45 * time.Second)).To(Equal(Seconds))
			Expect(FindUnit(65 * time.Second)).To(Equal(Minutes))
			Expect(FindUnit(time.Minute)).To(Equal(Minutes))
			Expect(FindUnit(25 * time.Minute)).To(Equal(Minutes))
			Expect(FindUnit(10 * time.Hour)).To(Equal(Hours))
			Expect(FindUnit(24 * time.Hour)).To(Equal(Days))
			Expect(FindUnit(7 * 24 * time.Hour)).To(Equal(Weeks))
			Expect(FindUnit(30 * 24 * time.Hour)).To(Equal(Months))
			Expect(FindUnit(365 * 24 * time.Hour)).To(Equal(Years))
		})
	})

})
