package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/psewda/typing"
)

var verFlag bool

func init() {
	parseFlags()
}

func main() {
	if verFlag {
		verStr := buildVerStr(typing.Version, typing.BuildNumber)
		fmt.Println(verStr)
		return
	}
	fmt.Println("Hello Typing !!")
}

func parseFlags() {
	flag.BoolVar(&verFlag, "version", false, "print typing version")
	flag.Parse()
}

func buildVerStr(ver, build string) string {
	osArch := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	return fmt.Sprintf("Typing %s-%s %s", ver, build, osArch)
}
