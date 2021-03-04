package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/psewda/typing"
	"github.com/psewda/typing/pkg/controllers"
	"github.com/psewda/typing/pkg/log"
	"github.com/psewda/typing/pkg/server"
)

const (
	envVarPort       = "TYPING_PORT"
	envVarLogLevel   = "TYPING_LOG_LEVEL"
	buildTypeDebug   = "DEBUG"
	buildTypeRelease = "RELEASE"
)

var (
	build   string = buildTypeDebug
	port    uint16
	logger  *log.Logger
	verFlag bool
)

func init() {
	// initialize cli args
	parseFlags()
	if verFlag {
		return
	}

	// initialize the logger
	logLevel, ok := parseLogLevel(os.Getenv(envVarLogLevel))
	if !ok {
		logLevel = log.LevelTypeDebug
	}
	config := log.Configuration{
		Level:  logLevel,
		Output: os.Stdout,
		Color:  true,
	}
	if build == buildTypeRelease {
		config.Color = false
	}
	logger = log.New(config)

	// set port for server
	p, ok := parsePort(os.Getenv(envVarPort))
	if !ok {
		p = server.GetRandPort()
	}
	port = p
}

func main() {
	// if --version is passed, print version string
	if verFlag {
		verStr := typing.GetVersionString()
		fmt.Println(verStr)
		return
	}

	// create new api server
	server := server.New(true, logger)

	// register api controllers
	server.RegisterController(controllers.NewVersionController())

	// run the api server
	if err := server.Run(port); err != nil {
		logger.Fatal("error occurred while starting the server", err)
	}
	logger.Info(fmt.Sprintf("server started on port %d, happy to serve api !", port))

	// wait for intrupt signal to exit
	waitForInterrupt(server)
	logger.Info("process terminated gracefully, have a wonderful day !")
}

func waitForInterrupt(server *server.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	logger.Info("caught interrupt signal, terminating process")
	server.Shutdown()
}

func parseFlags() {
	flag.BoolVar(&verFlag, "version", false, "print typing version")
	flag.Parse()
}

func parsePort(p string) (uint16, bool) {
	if len(p) > 0 {
		if v, err := strconv.Atoi(p); err == nil {
			if v >= 1024 && v <= 65535 {
				return uint16(v), true
			}
		}
	}
	return 0, false
}

func parseLogLevel(level string) (log.LevelType, bool) {
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case "DEBUG":
		return log.LevelTypeDebug, true
	case "INFO":
		return log.LevelTypeInfo, true
	case "WARN":
		return log.LevelTypeWarn, true
	case "ERROR":
		return log.LevelTypeError, true
	default:
		return 255, false
	}
}
