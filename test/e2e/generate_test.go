package e2e_test

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Generate", func() {
	It("should work", func() {
		outpath := GinkgoT().TempDir()
		cmd := exec.Command(cmdPath,
			"--package-name", "petstore",
			"--specification", petstoreSpecPath,
			"--config", filepath.Join(gitRoot, "test", "e2e", "testdata", "petstore", "openapi2go.yml"),
			"--output", outpath,
		)
		data, err := fs.ReadFile(testdata, "testdata/petstore/petstore.go")
		Expect(err).NotTo(HaveOccurred())
		expected := string(data)

		ses, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred())
		Eventually(ses).Should(gexec.Exit(0))

		genpath := filepath.Join(outpath, "petstore.go")
		Expect(genpath).To(BeARegularFile())
		actual, err := os.ReadFile(genpath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(actual)).To(Equal(expected))
	})
})
