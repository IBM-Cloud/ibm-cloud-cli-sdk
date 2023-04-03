package bluemix

import "fmt"

// Version is the SDK version
var Version = VersionType{Major: 1, Minor: 0, Build: 3}

// VersionType describe version info
type VersionType struct {
	Major int // major version
	Minor int // minor version
	Build int // build number
}

// String will return the version in semver format string "Major.Minor.Build"
func (v VersionType) String() string {
	if v == (VersionType{}) {
		return ""
	}
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Build)
}
