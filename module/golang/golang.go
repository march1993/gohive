package golang

import (
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/config"
	"github.com/march1993/gohive/module"
	_ "github.com/march1993/gohive/module/git"
	"github.com/march1993/gohive/module/linux"
	"github.com/march1993/gohive/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type golang struct{}

func init() {
	module.RegisterModule("golang", &golang{})
	if err := os.MkdirAll(config.GOLANG_DIR, config.GOLANG_DIR_PERM); err != nil {
		panic(err.Error())
	}
	if err := os.Chmod(config.GOLANG_DIR, config.GOLANG_DIR_PERM); err != nil {
		panic(err.Error())
	}
	if err := os.Chown(config.GOLANG_DIR, config.GOLANG_DIR_O_USER, config.GOLANG_DIR_O_GROUP); err != nil {
		panic(err.Error())
	}

}

func (g *golang) Create(name string) api.Status {
	// do nothing
	return api.Status{Status: api.STATUS_SUCCESS}
}

func (g *golang) Rename(oldName string, newName string) api.Status {
	// do nothing
	return api.Status{Status: api.STATUS_SUCCESS}
}

func (g *golang) Remove(name string) api.Status {
	// do nothing
	return api.Status{Status: api.STATUS_SUCCESS}
}

func (g *golang) Status(name string) api.Status {
	unixname := linux.Prefix + name
	stdout, err := exec.Command("runuser",
		unixname,
		"-s", "/bin/bash",
		"-c", "go version",
	).CombinedOutput()
	if err != nil {
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

func (g *golang) Repair(name string) api.Status {
	unixname := linux.Prefix + name
	errs := []string{}

	if stdout, err := exec.Command("runuser",
		unixname,
		"-s", "/bin/bash",
		"-c", "cd ~ && go get .",
	).CombinedOutput(); err != nil {
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

func (g *golang) ListRemoved() []string {
	return []string{}
}

const (
	GO_PREFIX        = "go"
	GO_DOWNLOAD_PATH = "https://dl.google.com/go/go{{VERSION}}.linux-amd64.tar.gz"
)

func GetGolangList() api.Status {

	result := []string{}

	files, err := ioutil.ReadDir(config.GOLANG_DIR)

	if err != nil {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: err.Error(),
		}
	}

	for _, file := range files {
		if file.IsDir() {
			name := file.Name()

			if !strings.HasPrefix(name, GO_PREFIX) {
				result = append(result, name)
			}
		}
	}

	return api.Status{
		Status: api.STATUS_SUCCESS,
		Result: result,
	}
}

func SetGolangInstallation(version string) api.Status {

	tmpPath := config.GOLANG_DIR + "/go" + version + ".tmp"
	tmp, err := os.Create(tmpPath)

	if err != nil {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: err.Error(),
		}
	}

	defer tmp.Close()
	defer os.Remove(tmpPath)

	resp, err := http.Get(strings.Replace(GO_DOWNLOAD_PATH, "{{VERSION}}", version, -1))

	if err != nil {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: err.Error(),
		}
	}

	defer resp.Body.Close()

	_, err = io.Copy(tmp, resp.Body)

	if err != nil {

		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: err.Error(),
		}
	}

	outputPath := config.GOLANG_DIR + "/go" + version

	if t := util.Shell("mkdir", outputPath); t != "" {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: t,
		}
	}

	if t := util.Shell("tar",
		"xf", tmpPath,
		"-C", outputPath,
		"--strip-components", "!"); t != "" {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: t,
		}
	}

	return api.Status{Status: api.STATUS_SUCCESS}
}
