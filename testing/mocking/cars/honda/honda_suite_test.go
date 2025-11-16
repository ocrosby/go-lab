package honda

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHonda(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Honda Suite")
}
