package e2e_test

import (
	"context"
	"embed"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/unmango/go/vcs/git"
)

var (
	cmdPath          string
	gitRoot          string
	petstoreSpecPath string

	//go:embed testdata
	testdata embed.FS
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2e Suite")
}

var _ = BeforeSuite(func(ctx context.Context) {
	root, err := git.Root(ctx)
	Expect(err).NotTo(HaveOccurred())
	gitRoot = root

	cmdPath, err = gexec.Build(root)
	Expect(err).NotTo(HaveOccurred())

	petstoreSpecPath = filepath.Join(root, "bin", "petstore.json")
	Expect(petstoreSpecPath).To(BeARegularFile())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
