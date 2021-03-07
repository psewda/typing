package utils_test

import (
	"errors"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/psewda/typing/internal/utils"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "utils-suite")
}

var _ = Describe("utility functions", func() {
	Context("utility function: GetValueString()", func() {
		It("should cover all cases", func() {
			By("valid value")
			v := utils.GetValueString("value", utils.Empty)
			Expect(v).Should(Equal("value"))

			By("empty value")
			v = utils.GetValueString(utils.Empty, "default")
			Expect(v).Should(Equal("default"))
		})
	})

	Context("utility function: Error()", func() {
		It("should cover all cases", func() {
			By("valid input")
			err := utils.Error("value", errors.New("inner"))
			Expect(err.Error()).Should(Equal("value: [inner]"))

			By("zero input")
			err = utils.Error(utils.Empty, nil)
			Expect(err.Error()).Should(Equal("error"))
		})
	})

	Context("utility function: AppendError()", func() {
		It("should cover all cases", func() {
			By("valid input")
			msg := utils.AppendError("value", errors.New("inner"))
			Expect(msg).Should(Equal("value: [inner]"))

			By("only value")
			msg = utils.AppendError("value", nil)
			Expect(msg).Should(Equal("value"))
		})
	})

	Context("utility function: CheckLocalhostURL()", func() {
		It("should cover all cases", func() {
			By("valid url")
			err := utils.CheckLocalhostURL("http://localhost:5050/redirect")
			Expect(err).ShouldNot(HaveOccurred())

			By("wrong url")
			err = utils.CheckLocalhostURL("invalid-value")
			Expect(err).Should(HaveOccurred())

			By("non localhost url")
			err = utils.CheckLocalhostURL("http://example.com:5050/redirect")
			Expect(err).Should(HaveOccurred())
		})
	})
})
