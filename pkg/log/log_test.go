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
			logger := newLogger(buffer, log.LevelTypeDebug, false)
			logger.Debug("debug")

			Expect(buffer.String()).Should(HavePrefix("{"))
			Expect(buffer.String()).Should(HaveSuffix("}\n"))
			Expect(strings.Split(buffer.String(), ",")).Should(HaveLen(3))
		})
	})

	Context("log content", func() {
		It("should log correct content", func() {
			buffer := new(bytes.Buffer)
			logger := newLogger(buffer, log.LevelTypeDebug, false)

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
		It("should log correct level - debug", func() {
			buffer := new(bytes.Buffer)
			logger := newLogger(buffer, log.LevelTypeDebug, false)
			writeAll(logger)
			Expect(buffer.String()).Should(ContainSubstring(`"level":"DEBUG"`))
			Expect(buffer.String()).Should(ContainSubstring(`"level":"INFO"`))
			Expect(buffer.String()).Should(ContainSubstring(`"level":"WARN"`))
			Expect(buffer.String()).Should(ContainSubstring(`"level":"ERROR"`))
		})

		It("should log correct level - info", func() {
			buffer := new(bytes.Buffer)
			logger := newLogger(buffer, log.LevelTypeInfo, false)
			writeAll(logger)
			Expect(buffer.String()).ShouldNot(ContainSubstring(`"level":"DEBUG"`))
			Expect(buffer.String()).Should(ContainSubstring(`"level":"INFO"`))
			Expect(buffer.String()).Should(ContainSubstring(`"level":"WARN"`))
			Expect(buffer.String()).Should(ContainSubstring(`"level":"ERROR"`))
		})

		It("should log correct level - warn", func() {
			buffer := new(bytes.Buffer)
			logger := newLogger(buffer, log.LevelTypeWarn, false)
			writeAll(logger)
			Expect(buffer.String()).ShouldNot(ContainSubstring(`"level":"DEBUG"`))
			Expect(buffer.String()).ShouldNot(ContainSubstring(`"level":"INFO"`))
			Expect(buffer.String()).Should(ContainSubstring(`"level":"WARN"`))
			Expect(buffer.String()).Should(ContainSubstring(`"level":"ERROR"`))
		})

		It("should log correct level - error", func() {
			buffer := new(bytes.Buffer)
			logger := newLogger(buffer, log.LevelTypeError, false)
			writeAll(logger)
			Expect(buffer.String()).ShouldNot(ContainSubstring(`"level":"DEBUG"`))
			Expect(buffer.String()).ShouldNot(ContainSubstring(`"level":"INFO"`))
			Expect(buffer.String()).ShouldNot(ContainSubstring(`"level":"WARN"`))
			Expect(buffer.String()).Should(ContainSubstring(`"level":"ERROR"`))
		})
	})

	Context("log color", func() {
		const colorValue = "["

		It("should enable color", func() {
			buffer := new(bytes.Buffer)
			logger := newLogger(buffer, log.LevelTypeDebug, true)
			logger.Debug("debug")
			Expect(buffer.String()).Should(ContainSubstring(colorValue))
		})

		It("should disable color", func() {
			buffer := new(bytes.Buffer)
			logger := newLogger(buffer, log.LevelTypeDebug, false)
			logger.Debug("debug")
			Expect(buffer.String()).ShouldNot(ContainSubstring(colorValue))
		})
	})
})

func newLogger(w io.Writer, level log.LevelType, color bool) *log.Logger {
	return log.New(log.Configuration{
		Output: w,
		Level:  level,
		Color:  color,
	})
}

func writeAll(logger *log.Logger) {
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error", errors.New("error"))
}
