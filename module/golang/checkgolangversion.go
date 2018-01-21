package golang

import (
	"regexp"
)

func CheckGolangVersion(version string) bool {
	re := regexp.MustCompile("^[0-9.]{1,16}$")
	return re.MatchString(version)
}
