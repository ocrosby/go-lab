package toyota

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestToyota(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Toyota Suite")
}
