package openapi2go_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOpenapi2Go(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Openapi2Go Suite")
}
