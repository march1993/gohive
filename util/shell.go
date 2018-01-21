package util

import (
	"os/exec"
	. "strings"
)

// warning: running as root
func Shell(name string, arg ...string) string {
	bytes, _ := exec.Command(name, arg...).CombinedOutput()
	output := string(bytes)
	output = Trim(output, "\n")
	return string(output)
}
