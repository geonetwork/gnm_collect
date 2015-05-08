package tests
import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPrivateGnmsysAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gnmsys Private API Test Suite")
}
