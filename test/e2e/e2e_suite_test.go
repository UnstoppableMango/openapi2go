package e2e_test

import (
	"context"
	"embed"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	cmdPath          string
	petstoreSpecPath string

	//go:embed testdata
	testdata embed.FS
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2e Suite", Label("E2E"))
}

var _ = BeforeSuite(func(ctx context.Context) {
	var err error
	cmdPath, err = gexec.Build("../..")
	Expect(err).NotTo(HaveOccurred())

	petstoreSpecPath = filepath.Join("../..", "bin", "petstore.json")
	Expect(petstoreSpecPath).To(BeARegularFile())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
