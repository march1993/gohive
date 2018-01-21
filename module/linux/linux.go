package linux

import (
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/config"
	"github.com/march1993/gohive/module"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

type linux struct{}

const (
	Prefix = "gohive_app_"
	Group  = "gohive_app"
	Suffix = ".data"
)

func init() {
	module.RegisterModule("linux", &linux{})

	cmd := exec.Command("groupadd", "-f", Group)
	if stdout, err := cmd.CombinedOutput(); err != nil {
		panic(string(stdout) + err.Error())
	}

}

func getHomeDir(name string) string {
	return config.APP_DIR + "/" + Prefix + name
}

func getDataDir(name string) string {
	return config.APP_DIR + "/" + Prefix + name + Suffix
}

func (l *linux) Create(name string) api.Status {

	if l.Status(name).Reason == api.APP_NON_EXIST {
		unixname := Prefix + name

		cmd := exec.Command("useradd",
			"-b", config.APP_DIR, // home directory
			"-m",                   // create home
			"-s", config.SSH_SHELL, // shell
			"-g", Group, // group
			"-K", "UMASK=0077",
			unixname)
		stdout, err := cmd.CombinedOutput()

		os.MkdirAll(getDataDir(name), 0700)

		if err != nil {
			return api.Status{
				Status: api.STATUS_FAILURE,
				Reason: string(stdout),
			}
		}

		return api.Status{Status: api.STATUS_SUCCESS}

	} else {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.APP_ALREADY_EXISTING,
		}
	}
}

func (l *linux) Rename(oldName string, newName string) api.Status {
	if l.Status(oldName).Status != api.STATUS_SUCCESS {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.REASON_CONDITION_UNMET,
		}
	} else {
		// TODO:
		// 1. kill all process
		// 2. rename user
		// 3. rename folders
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.REASON_UNKNOWN,
		}
	}
}

func (l *linux) Remove(name string) api.Status {
	if ret := l.Status(name); ret.Reason == api.APP_NON_EXIST {
		return ret
	} else {
		unixname := Prefix + name
		cmd := exec.Command("userdel", unixname)
		cmd.CombinedOutput()

		os.RemoveAll(getHomeDir(name))
		os.RemoveAll(getDataDir(name))

		return api.Status{Status: api.STATUS_SUCCESS}
	}
}

func (l *linux) Status(name string) api.Status {

	unixname := Prefix + name

	parital := false
	broken := false

	// check user
	_, err := user.Lookup(unixname)
	if err != nil {
		broken = true
	} else {
		parital = true
	}

	// check group
	cmd := exec.Command("id", "-gn", unixname)
	stdout, err := cmd.CombinedOutput()
	if err != nil || strings.Trim(string(stdout), "\n") != Group {
		broken = true
	} else {
		parital = true
	}

	// check files
	cmd = exec.Command("stat", "-c", "%U", getHomeDir(name))
	stdout, err = cmd.CombinedOutput()
	if err != nil || strings.Trim(string(stdout), "\n") != unixname {
		broken = true
	} else {
		parital = true
	}

	cmd = exec.Command("stat", "-c", "%G", getHomeDir(name))
	stdout, err = cmd.CombinedOutput()
	if err != nil || strings.Trim(string(stdout), "\n") != Group {
		broken = true
	} else {
		parital = true
	}

	cmd = exec.Command("stat", "-c", "%U", getDataDir(name))
	stdout, err = cmd.CombinedOutput()
	if err != nil || strings.Trim(string(stdout), "\n") != "root" {
		broken = true
	} else {
		parital = true
	}

	cmd = exec.Command("stat", "-c", "%G", getDataDir(name))
	stdout, err = cmd.CombinedOutput()
	if err != nil || strings.Trim(string(stdout), "\n") != "root" {
		broken = true
	} else {
		parital = true
	}

	if broken {
		if parital {
			return api.Status{
				Status: api.STATUS_FAILURE,
				Reason: api.APP_BROKEN,
			}
		} else {
			return api.Status{
				Status: api.STATUS_FAILURE,
				Reason: api.APP_NON_EXIST,
			}
		}
	} else {
		return api.Status{Status: api.STATUS_SUCCESS}
	}

}

func (l *linux) Repair(name string) api.Status {

	// create user
	_, err := user.Lookup(name)
	if err != nil {
		l.Create(name)
	}

	// fix files owners
	unixname := Prefix + name
	cmd := exec.Command("chown",
		unixname+":"+Group,
		"-R",
		getHomeDir(name))
	stdout, err := cmd.CombinedOutput()

	if err != nil {
		panic(string(stdout) + err.Error())
	}

	cmd = exec.Command("chown",
		"root:root",
		"-R",
		getDataDir(name))
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		panic(string(stdout) + err.Error())
	}

	// fix group
	cmd = exec.Command("usermod",
		"-g", Group,
		unixname)
	stdout, err = cmd.CombinedOutput()

	if err != nil {
		panic(string(stdout) + err.Error())
	}

	return api.Status{Status: api.STATUS_SUCCESS}
}

func (l *linux) ListRemoved() []string {
	cmd := exec.Command("members", Group)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	members := strings.Split(strings.Trim(string(stdout), "\n"), " ")

	ret := []string{}

	for _, member := range members {
		if l.Status(member).Status != api.STATUS_SUCCESS {
			ret = append(ret, member)
		}
	}

	return ret
}
