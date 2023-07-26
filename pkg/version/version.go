package version

import (
	"fmt"
	"runtime"
)

var (
	Version = "0.1.0"
	// Git SHA Value will be set during build
	gitCommit = "Not provided (use make build instead of go build)"
	buildDate = "1970-01-01T00:00:00Z" // build date, output of $(date +'%Y-%m-%dT%H:%M:%S')
)

func Print() {
	fmt.Printf(`Name: harbor-cleaner
Version: %s
CommitID: %s
BuildDate: %s
GoVersion: %s
Compiler: %s
Platform: %s/%s
`, Version, gitCommit, buildDate, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
}
