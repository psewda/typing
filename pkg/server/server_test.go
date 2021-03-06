package server_test

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"github.com/psewda/typing/pkg/log"
	"github.com/psewda/typing/pkg/server"
)

func TestServer(t *testing.T) {
	DefaultReporterConfig.SlowSpecThreshold = 30
	RegisterFailHandler(Fail)
	RunSpecs(t, "server-suite")
}

var _ = Describe("server", func() {
	Context("run server", func() {
		It("should succeed when valid port", func() {
			s := server.New(false, newLogger())
			err := s.Run(server.GetRandPort())
			defer s.Shutdown()
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(s.Running).Should(BeTrue())
		})

		It("should return error when invalid port", func() {
			s := server.New(false, newLogger())
			err := s.Run(25)
			Expect(err).Should(HaveOccurred())
			Expect(s.Running).Should(BeFalse())
		})

		It("should return error when non-free port", func() {
			s := server.New(false, newLogger())
			s.Run(5500)
			defer s.Shutdown()
			s2 := server.New(false, newLogger())
			err := s2.Run(5500)
			Expect(err).Should(HaveOccurred())
			Expect(s2.Running).Should(BeFalse())
		})
	})

	Context("shutdown server", func() {
		It("should succeed when valid setup", func() {
			s := server.New(false, newLogger())
			err := s.Run(server.GetRandPort())
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(s.Running).Should(BeTrue())

			s.Shutdown()
			Expect(s.Running).Should(BeFalse())
		})
	})
})

func newLogger() *log.Logger {
	return log.New(log.Configuration{
		Output: ioutil.Discard,
	})
}
