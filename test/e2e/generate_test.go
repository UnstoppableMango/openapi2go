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
		outdir := GinkgoT().TempDir()
		cmd := exec.Command(cmdPath, petstoreSpecPath,
			"--package-name", "petstore",
			"--output", outdir,
		)
		data, err := fs.ReadFile(testdata, "testdata/petstore.go")
		Expect(err).NotTo(HaveOccurred())
		expected := string(data)

		ses, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred())
		Eventually(ses).Should(gexec.Exit(0))

		genpath := filepath.Join(outdir, "petstore.go")
		Expect(genpath).To(BeARegularFile())
		actual, err := os.ReadFile(genpath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(actual)).To(Equal(expected))
	})
})
