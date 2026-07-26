package mocking_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMocking(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mocking Suite")
}
