package infraestructure

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConsumptiServices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Consumption Handler Test Suite")
}
