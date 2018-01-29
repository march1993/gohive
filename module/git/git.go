package git

import (
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/config"
	"github.com/march1993/gohive/module"
	_ "github.com/march1993/gohive/module/linux"
	"github.com/march1993/gohive/module/systemd"
	"github.com/march1993/gohive/util"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const (
	gitPostUpdateTemplate  = "./templates/post-update"
	gitPostUpdateGenerated = "/repo.git/hooks/post-update"
)

type git struct{}

func init() {
	module.RegisterModule("git", &git{})
}

func (g *git) Create(name string) api.Status {
	unixname := config.APP_PREFIX + name
	errs := []string{}

	/* 1. create repo.git */
	stdout, err := exec.Command("runuser",
		unixname,
		"-s", "/bin/bash",
		"-c", "cd ~ && git init --bare repo.git",
	).CombinedOutput()
	if err != nil {
		errs = append(errs, string(stdout))

	}

	/* 2. copy post-update hook file */
	generated := config.GetHomeDir(name) + gitPostUpdateGenerated
	bytes, err := ioutil.ReadFile(gitPostUpdateTemplate)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		template := string(bytes)
		template = strings.Replace(template, "{{GOLANG_EXECUTABLE}}", config.GOLANG_EXECUTABLE, -1)
		template = strings.Replace(template, "{{UNIXNAME}}", unixname, -1)
		template = strings.Replace(template, "{{GOHIVE}}", systemd.ExecStart, -1)
		config.AppConfigSet(name, "git", "post-update-hash", util.Hash(template))

		err = ioutil.WriteFile(generated, []byte(template), 0755)

		if err != nil {
			errs = append(errs, err.Error())
		}
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

func (g *git) Rename(oldName string, newName string) api.Status {
	// do nothing
	return api.Status{Status: api.STATUS_SUCCESS}
}

func (g *git) Remove(name string) api.Status {
	// do nothing
	return api.Status{Status: api.STATUS_SUCCESS}
}

func (g *git) Status(name string) api.Status {
	unixname := config.APP_PREFIX + name
	errs := []string{}

	stdout, err := exec.Command("runuser",
		unixname,
		"-s", "/bin/bash",
		"-c", "cd ~/repo.git && git branch -a",
	).CombinedOutput()
	if err != nil {
		errs = append(errs, string(stdout))
	}

	bytes, _ := ioutil.ReadFile(config.GetHomeDir(name) + gitPostUpdateGenerated)
	hash := config.AppConfigGet(name, "git", "post-update-hash", "")

	if err != nil {
		errs = append(errs, err.Error())
	} else {

		if hash != util.Hash(string(bytes)) || "" == hash {
			errs = append(errs, api.GIT_POST_UPDATE_EXPIRED)
		}
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

func (g *git) Repair(name string) api.Status {
	return g.Create(name)
}

func (g *git) ListRemoved() []string {
	return []string{}
}

func GetGitUrl(name string) string {
	unixname := config.APP_PREFIX + name
	return unixname + "@" + config.Get("server_name", "${server_name}") + ":~/repo.git"
}

const (
	SSH_DIR      = "/.ssh"
	SSH_KEY_FILE = SSH_DIR + "/authorized_keys"
	SSH_KEY_PERM = 0644
)

func SetGitKeys(name string, keys []string) api.Status {
	unixname := config.APP_PREFIX + name
	home := config.GetHomeDir(name)

	errs := []string{}

	if err := os.MkdirAll(home+SSH_DIR, 0700); err != nil {
		errs = append(errs, err.Error())
	}

	if err := ioutil.WriteFile(home+SSH_KEY_FILE, []byte(strings.Join(keys, "\n")), SSH_KEY_PERM); err != nil {
		errs = append(errs, err.Error())
	}

	if stdout, err := exec.Command("chown",
		unixname+":"+config.APP_GROUP,
		"-R",
		home+SSH_DIR).CombinedOutput(); err != nil {
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

func GetGitKeys(name string) api.Status {
	home := config.GetHomeDir(name)

	if stdout, err := exec.Command("ssh-keygen",
		"-lf",
		home+SSH_KEY_FILE).CombinedOutput(); err != nil {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: string(stdout),
		}
	} else {
		return api.Status{
			Status: api.STATUS_SUCCESS,
			Result: string(stdout),
		}
	}
}
