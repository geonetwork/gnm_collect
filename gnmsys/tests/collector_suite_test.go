package tests
import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAddPoint(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gnmsys Test Suite")
}
