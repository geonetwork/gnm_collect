package gnmsys
import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPrivateGnmSysApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gnmsys Private API Test Suite")
}
