package linux

import (
	"errors"
	"github.com/march1993/gohive/api"
	. "github.com/march1993/gohive/config"
	"github.com/march1993/gohive/module"
	"github.com/march1993/gohive/util"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

type linux struct{}

func init() {
	module.RegisterModule("linux", &linux{})

	cmd := exec.Command("groupadd", "-f", APP_GROUP)
	if stdout, err := cmd.CombinedOutput(); err != nil {
		panic(string(stdout) + err.Error())
	}

}

func (l *linux) Create(name string) api.Status {

	if l.Status(name).Reason == api.APP_NON_EXIST {
		unixname := APP_PREFIX + name

		cmd := exec.Command("useradd",
			"-b", APP_DIR, // home directory
			"-m",            // create home
			"-s", SSH_SHELL, // shell
			"-g", APP_GROUP, // group
			"-K", "UMASK=0077",
			unixname)
		stdout, err := cmd.CombinedOutput()

		os.MkdirAll(GetDataDir(name), 0700)

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

func (l *linux) Rename(oldName, newName string) api.Status {
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
		if stdout, err := exec.Command("killall", "--user", APP_PREFIX+oldName).CombinedOutput(); err != nil {
			if string(stdout) != "" {
				errs = append(errs, string(stdout))
			}
		}

		if stdout, err := exec.Command("killall", "-s", "9", "--user", APP_PREFIX+oldName).CombinedOutput(); err != nil {
			if string(stdout) != "" {
				errs = append(errs, string(stdout))
			}
		}

		// 2. rename user
		if stdout, err := exec.Command("usermod",
			"-l", APP_PREFIX+newName,
			"-m", "-d", GetHomeDir(newName),
			APP_PREFIX+oldName).CombinedOutput(); err != nil {
			errs = append(errs, string(stdout))
		}

		// 3. rename folders
		if err := os.Rename(GetDataDir(oldName), GetDataDir(newName)); err != nil {
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
		unixname := APP_PREFIX + name

		errs := []string{}
		// 1. kill all process
		if stdout, err := exec.Command("killall", "--user", APP_PREFIX+name).CombinedOutput(); err != nil {
			if string(stdout) != "" {
				errs = append(errs, string(stdout))
			}
		}
		if stdout, err := exec.Command("killall", "-s", "9", "--user", APP_PREFIX+name).CombinedOutput(); err != nil {
			if string(stdout) != "" {
				errs = append(errs, string(stdout))
			}
		}

		// 2. remove the user
		if stdout, err := exec.Command("userdel", unixname).CombinedOutput(); err != nil {
			errs = append(errs, string(stdout))
		}

		// 3. remove the directories
		if err := os.RemoveAll(GetHomeDir(name)); err != nil {
			errs = append(errs, err.Error())
		}
		if err := os.RemoveAll(GetDataDir(name)); err != nil {
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

	unixname := APP_PREFIX + name

	parital := false
	errs := []string{}

	// check user
	if _, err := user.Lookup(unixname); err != nil {
		errs = append(errs, err.Error())
	} else {
		parital = true
	}

	// check group
	if err := testCmdHelper(APP_GROUP, "id", "-gn", unixname); err != nil {
		errs = append(errs, err.Error())
	} else {
		parital = true
	}

	// check files
	if err := testCmdHelper(unixname, "stat", "-c", "%U", GetHomeDir(name)); err != nil {
		errs = append(errs, err.Error())
	} else {
		parital = true
	}

	if err := testCmdHelper(APP_GROUP, "stat", "-c", "%G", GetHomeDir(name)); err != nil {
		errs = append(errs, err.Error())
	} else {
		parital = true
	}

	if err := testCmdHelper("root", "stat", "-c", "%U", GetDataDir(name)); err != nil {
		errs = append(errs, err.Error())
	} else {
		parital = true
	}

	if err := testCmdHelper("root", "stat", "-c", "%G", GetDataDir(name)); err != nil {
		errs = append(errs, err.Error())
	} else {
		parital = true
	}

	// check chmod
	if err := testCmdHelper("700", "stat", "-c", "%a", GetHomeDir(name)); err != nil {
		errs = append(errs, err.Error())
	} else {
		parital = true
	}

	if err := testCmdHelper("700", "stat", "-c", "%a", GetDataDir(name)); err != nil {
		errs = append(errs, err.Error())
	} else {
		parital = true
	}

	if len(errs) > 0 {
		if parital {
			return api.Status{
				Status: api.STATUS_FAILURE,
				Reason: api.APP_BROKEN,
				Result: strings.Join(errs, "\n"),
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
	// ensure directories exist
	if err := os.MkdirAll(GetHomeDir(name), 0700); err != nil {
		errs = append(errs, err.Error())
	}
	if err := os.MkdirAll(GetDataDir(name), 0700); err != nil {
		errs = append(errs, err.Error())
	}
	// fix files owners
	unixname := APP_PREFIX + name
	stdout, err := exec.Command("chown",
		unixname+":"+APP_GROUP,
		"-R",
		GetHomeDir(name)).CombinedOutput()
	if err != nil {
		errs = append(errs, string(stdout))
	}

	stdout, err = exec.Command("chown",
		"root:root",
		"-R",
		GetDataDir(name)).CombinedOutput()
	if err != nil {
		errs = append(errs, string(stdout))
	}

	// chmod
	err = os.Chmod(GetHomeDir(name), 0700)
	if err != nil {
		errs = append(errs, err.Error())
	}

	err = os.Chmod(GetDataDir(name), 0700)
	if err != nil {
		errs = append(errs, err.Error())
	}

	// fix group and shell
	stdout, err = exec.Command("usermod",
		"-s", SSH_SHELL,
		"-g", APP_GROUP,
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
	// 1. list remaining users
	cmd := exec.Command("members", APP_GROUP)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		panic(string(stdout) + err.Error())
	}

	members := strings.Split(strings.Trim(string(stdout), "\n"), " ")

	ret := []string{}
	list := GetAppList()

	for _, member := range members {
		if !util.Includes(list, member) {
			ret = append(ret, member)
		}
	}

	// 2. TODO: list remaining process whose owner has been deleted

	return ret
}

func GetAppList() []string {
	result := []string{}

	files, err := ioutil.ReadDir(APP_DIR)

	if err != nil {
		panic(err.Error())
	}

	for _, file := range files {
		if file.IsDir() {
			name := file.Name()

			if !strings.HasSuffix(name, APP_DATA_SUFFIX) {
				result = append(result, name)
			}
		}
	}

	return result
}
