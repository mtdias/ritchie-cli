package team

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRitTeam(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rit Suite")
}

var _ = Describe("RitScaffold", func() {
	BeforeSuite(func() {
	})
})
