package version

import (
    "strings"
)

var Version = "0.0.1"

func GetCurrentVersion(mod string) string {
    return Version
}

func GetMinorVersion(version string) string {
	versionList := strings.Split(version, ".")
	if len(versionList) < 3 {
		return ""
	}
	return versionList[0] + "." + versionList[1]
}
