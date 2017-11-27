package bluemix

var Version = VersionType{Major: 1, Minor: 0, Build: 0}

type VersionType struct {
	Major int
	Minor int
	Build int
}
