package bluemix

// Version is the SDK version
var Version = VersionType{Major: 0, Minor: 1, Build: 0}

// VersionType describe version info
type VersionType struct {
	Major int // major version
	Minor int // minor version
	Build int // build number
}
