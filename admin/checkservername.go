package admin

import (
	"regexp"
	. "strings"
)

func checkServerName(name string) bool {

	re := regexp.MustCompile("^[0-9a-zA-Z_-.]{1,32}$")
	if HasPrefix(name, "-") {
		return false
	} else {
		return re.MatchString(name)
	}

}
