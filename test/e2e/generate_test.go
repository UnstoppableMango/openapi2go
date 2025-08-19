package e2e_test

import (
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Generate", func() {
	It("should work", func() {
		outpath := GinkgoT().TempDir()
		cmd := exec.Command(cmdPath, "generate",
			"--package-name", "petstore",
			"--specification", petstorePath,
			"--output", outpath,
		)

		ses, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred())
		Eventually(ses).Should(gexec.Exit(0))

		orderpath := filepath.Join(outpath, "petstore.go")
		Expect(orderpath).To(BeARegularFile())
	})
})
