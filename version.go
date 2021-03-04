package typing

import (
	"fmt"
	"runtime"
)

// Version is the place holder to put version string.
var Version = "0.1.0"

// BuildNumber is the place holder to put build-number
// string. It is set dynamically in the ci/cd pipeline.
var BuildNumber = "1"

// GetVersionString builds the version string.
func GetVersionString() string {
	osArch := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	return fmt.Sprintf("Typing %s-%s %s", Version, BuildNumber, osArch)
}
