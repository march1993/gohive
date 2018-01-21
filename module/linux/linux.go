package linux

import (
	"errors"
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
	} else if l.Status(newName).Reason != api.APP_NON_EXIST {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.APP_ALREADY_OCCUPIED,
		}
	} else {

		errs := []string{}
		// 1. kill all process
		if stdout, err := exec.Command("killall", "--user", Prefix+oldName).CombinedOutput(); err != nil {
			if string(stdout) != "" {
				errs = append(errs, string(stdout))
			}
		}

		if stdout, err := exec.Command("killall", "-s", "9", "--user", Prefix+oldName).CombinedOutput(); err != nil {
			if string(stdout) != "" {
				errs = append(errs, string(stdout))
			}
		}

		// 2. rename user
		if stdout, err := exec.Command("usermod",
			"-l", Prefix+newName,
			"-m", "-d", getHomeDir(newName),
			Prefix+oldName).CombinedOutput(); err != nil {
			errs = append(errs, string(stdout))
		}

		// 3. rename folders
		if err := os.Rename(getDataDir(oldName), getDataDir(newName)); err != nil {
			errs = append(errs, err.Error())
		}

		if len(errs) > 0 {
			return api.Status{
				Status: api.STATUS_FAILURE,
				Reason: strings.Join(errs, "\n"),
			}
		} else {
			return api.Status{Status: api.STATUS_SUCCESS}
		}
	}
}

func (l *linux) Remove(name string) api.Status {
	if ret := l.Status(name); ret.Reason == api.APP_NON_EXIST {
		return ret
	} else {
		unixname := Prefix + name

		errs := []string{}
		// 1. kill all process
		if stdout, err := exec.Command("killall", "--user", Prefix+name).CombinedOutput(); err != nil {
			if string(stdout) != "" {
				errs = append(errs, string(stdout))
			}
		}
		if stdout, err := exec.Command("killall", "-s", "9", "--user", Prefix+name).CombinedOutput(); err != nil {
			if string(stdout) != "" {
				errs = append(errs, string(stdout))
			}
		}

		// 2. remove the user
		if stdout, err := exec.Command("userdel", unixname).CombinedOutput(); err != nil {
			errs = append(errs, string(stdout))
		}

		// 3. remove the directories
		if err := os.RemoveAll(getHomeDir(name)); err != nil {
			errs = append(errs, err.Error())
		}
		if err := os.RemoveAll(getDataDir(name)); err != nil {
			errs = append(errs, err.Error())
		}

		if len(errs) > 0 {
			return api.Status{
				Status: api.STATUS_FAILURE,
				Reason: strings.Join(errs, "\n"),
			}
		} else {
			return api.Status{Status: api.STATUS_SUCCESS}
		}
	}
}

func testCmdHelper(desired string, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.CombinedOutput()
	trimmed := strings.Trim(string(stdout), "\n")
	if err != nil || trimmed != desired {
		return errors.New("[" + name + " " + strings.Join(arg, " ") + "]\n" +
			"Unexpected value: " + trimmed + " (" + desired + " desired)")
	} else {
		return nil
	}
}

func (l *linux) Status(name string) api.Status {

	unixname := Prefix + name

	parital := false
	errs := []string{}

	// check user
	if _, err := user.Lookup(unixname); err != nil {
		errs = append(errs, err.Error())
	} else {
		parital = true
	}

	// check group
	if err := testCmdHelper(Group, "id", "-gn", unixname); err != nil {
		errs = append(errs, err.Error())
	} else {
		parital = true
	}

	// check files
	if err := testCmdHelper(unixname, "stat", "-c", "%U", getHomeDir(name)); err != nil {
		errs = append(errs, err.Error())
	} else {
		parital = true
	}

	if err := testCmdHelper(Group, "stat", "-c", "%G", getHomeDir(name)); err != nil {
		errs = append(errs, err.Error())
	} else {
		parital = true
	}

	if err := testCmdHelper("root", "stat", "-c", "%U", getDataDir(name)); err != nil {
		errs = append(errs, err.Error())
	} else {
		parital = true
	}

	if err := testCmdHelper("root", "stat", "-c", "%G", getDataDir(name)); err != nil {
		errs = append(errs, err.Error())
	} else {
		parital = true
	}

	if len(errs) > 0 {
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

	errs := []string{}
	// fix files owners
	unixname := Prefix + name
	stdout, err := exec.Command("chown",
		unixname+":"+Group,
		"-R",
		getHomeDir(name)).CombinedOutput()
	if err != nil {
		errs = append(errs, string(stdout))
	}

	stdout, err = exec.Command("chown",
		"root:root",
		"-R",
		getDataDir(name)).CombinedOutput()
	if err != nil {
		errs = append(errs, string(stdout))
	}

	// fix group
	stdout, err = exec.Command("usermod",
		"-g", Group,
		unixname).CombinedOutput()
	if err != nil {
		errs = append(errs, string(stdout))
	}

	if len(errs) > 0 {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: strings.Join(errs, "\n"),
		}
	} else {
		return api.Status{Status: api.STATUS_SUCCESS}

	}
}

func (l *linux) ListRemoved() []string {
	cmd := exec.Command("members", Group)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		panic(string(stdout) + err.Error())
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
