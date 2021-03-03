package log_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/psewda/typing/pkg/log"
)

func TestLog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "log-suite")
}

var _ = Describe("logger", func() {
	Context("log format", func() {
		It("should log in correct format", func() {
			buffer := new(bytes.Buffer)
			logger := newLogger(buffer, log.LevelTypeDebug)

			By("json format")
			logger.Debug("debug")
			Expect(buffer.String()).Should(HavePrefix("{"))
			Expect(buffer.String()).Should(HaveSuffix("}\n"))

			By("only 3 fields")
			buffer.Reset()
			logger.Debug("debug")
			Expect(strings.Split(buffer.String(), ",")).Should(HaveLen(3))
		})
	})

	Context("log content", func() {
		It("should log correct content", func() {
			buffer := new(bytes.Buffer)
			logger := newLogger(buffer, log.LevelTypeDebug)

			By("content")
			logger.Debug("debug")
			logger.Info("info")
			logger.Warn("warn")
			logger.Error("error", errors.New("error"))
			Expect(buffer.String()).Should(ContainSubstring(`"message":"debug"`))
			Expect(buffer.String()).Should(ContainSubstring(`"message":"info"`))
			Expect(buffer.String()).Should(ContainSubstring(`"message":"warn"`))
			Expect(buffer.String()).Should(ContainSubstring(`"message":"error: [error]"`))
		})
	})

	Context("log level", func() {
		It("should log correct level", func() {
			buffer := new(bytes.Buffer)

			By("debug level")
			logger := newLogger(buffer, log.LevelTypeDebug)
			writeAll(logger)
			Expect(buffer.String()).Should(ContainSubstring(`"level":"DEBUG"`))
			Expect(buffer.String()).Should(ContainSubstring(`"level":"INFO"`))
			Expect(buffer.String()).Should(ContainSubstring(`"level":"WARN"`))
			Expect(buffer.String()).Should(ContainSubstring(`"level":"ERROR"`))

			By("info level")
			logger = newLogger(buffer, log.LevelTypeInfo)
			buffer.Reset()
			writeAll(logger)
			Expect(buffer.String()).ShouldNot(ContainSubstring(`"level":"DEBUG"`))
			Expect(buffer.String()).Should(ContainSubstring(`"level":"INFO"`))
			Expect(buffer.String()).Should(ContainSubstring(`"level":"WARN"`))
			Expect(buffer.String()).Should(ContainSubstring(`"level":"ERROR"`))

			By("warn level")
			logger = newLogger(buffer, log.LevelTypeWarn)
			buffer.Reset()
			writeAll(logger)
			Expect(buffer.String()).ShouldNot(ContainSubstring(`"level":"DEBUG"`))
			Expect(buffer.String()).ShouldNot(ContainSubstring(`"level":"INFO"`))
			Expect(buffer.String()).Should(ContainSubstring(`"level":"WARN"`))
			Expect(buffer.String()).Should(ContainSubstring(`"level":"ERROR"`))

			By("error level")
			logger = newLogger(buffer, log.LevelTypeError)
			buffer.Reset()
			writeAll(logger)
			Expect(buffer.String()).ShouldNot(ContainSubstring(`"level":"DEBUG"`))
			Expect(buffer.String()).ShouldNot(ContainSubstring(`"level":"INFO"`))
			Expect(buffer.String()).ShouldNot(ContainSubstring(`"level":"WARN"`))
			Expect(buffer.String()).Should(ContainSubstring(`"level":"ERROR"`))
		})
	})

})

func newLogger(w io.Writer, level log.LevelType) *log.Logger {
	return log.New(log.Configuration{
		Output: w,
		Color:  false,
		Level:  level,
	})
}

func writeAll(logger *log.Logger) {
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error", errors.New("error"))
}
