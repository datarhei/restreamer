package app

import (
	"fmt"
	"runtime"
)

// Name of the app
const Name = "datarhei-core"

type versionInfo struct {
	Major int
	Minor int
	Patch int
}

func (v versionInfo) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v versionInfo) MajorString() string {
	return fmt.Sprintf("%d.0.0", v.Major)
}

func (v versionInfo) MinorString() string {
	return fmt.Sprintf("%d.%d.0", v.Major, v.Minor)
}

// Version of the app
var Version = versionInfo{
	Major: 16,
	Minor: 16,
	Patch: 0,
}

// Commit is the git commit the app is build from. It should be filled in during compilation
var Commit = ""

// Branch is the git branch the app is build from. It should be filled in during compilation
var Branch = ""

// Build is the timestamp of when the app has been build. It should be filled in during compilation
var Build = ""

// Arch is the OS and CPU architecture this app is build for.
var Arch = runtime.GOOS + "/" + runtime.GOARCH

// Compiler is the golang version this app has been build with.
var Compiler = runtime.Version()
