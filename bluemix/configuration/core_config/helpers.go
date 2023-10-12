package core_config

import "runtime"

func DeterminePlatform() string {
	arch := runtime.GOARCH

	switch runtime.GOOS {
	case "windows":
		if arch == "386" {
			return "win32"
		} else {
			return "win64"
		}
	case "linux":
		switch arch {
		case "386":
			return "linux32"
		case "amd64":
			return "linux64"
		case "ppc64le":
			return "ppc64le"
		case "s390x":
			return "s390x"
		case "arm64":
			return "linux-arm64"
		}
	case "darwin":
		switch arch {
		case "arm64":
			return "osx-arm64"
		default:
			return "osx"
		}
	}
	return "unknown"
}
