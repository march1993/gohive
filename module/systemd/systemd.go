package systemd

import (
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/config"
	"github.com/march1993/gohive/module"
	"github.com/march1993/gohive/module/linux"
	"github.com/march1993/gohive/util"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type systemd struct{}

const (
	SYSTEMD_OUTPUT         = "/lib/systemd/system"
	SYSTEMD_TEMPLATE       = "./templates/app.service"
	SYSTEMD_SERVICE_SUFFIX = ".service"
)

func init() {
	module.RegisterModule("systemd", &systemd{})
}

func (s *systemd) Create(name string) api.Status {
	unixname := config.APP_PREFIX + name
	errs := []string{}

	bytes, err := ioutil.ReadFile(SYSTEMD_TEMPLATE)
	if err != nil {
		panic(err.Error())
	}

	generated := string(bytes)
	generated = strings.Replace(generated, "{{Description}}", "Systemd services for "+unixname, -1)
	generated = strings.Replace(generated, "{{ExecStart}}", config.GetHomeDir(name)+"/"+config.GOLANG_EXECUTABLE, -1)
	generated = strings.Replace(generated, "{{WorkingDirectory}}", config.GetHomeDir(name), -1)
	generated = strings.Replace(generated, "{{User}}", unixname, -1)
	generated = strings.Replace(generated, "{{Group}}", config.APP_GROUP, -1)
	generated = strings.Replace(generated, "{{Environment}}", getEnvironment(name), -1)

	hash := util.Hash(generated)
	config.AppConfigSet(name, "systemd", "hash", hash)

	output := SYSTEMD_OUTPUT + "/" + unixname + SYSTEMD_SERVICE_SUFFIX
	if err = ioutil.WriteFile(output, []byte(generated), 0644); err != nil {
		errs = append(errs, err.Error())
	}

	cmd := exec.Command("systemctl", "enable", unixname)
	if stdout, err := cmd.CombinedOutput(); err != nil {
		errs = append(errs, string(stdout))
	}

	cmd = exec.Command("systemctl", "restart", unixname)
	if stdout, err := cmd.CombinedOutput(); err != nil {
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

func (s *systemd) Rename(oldName, newName string) api.Status {

	if status := s.Remove(oldName); status.Status != api.STATUS_SUCCESS {
		return status
	}
	if status := s.Create(newName); status.Status != api.STATUS_SUCCESS {
		return status
	}

	return api.Status{Status: api.STATUS_SUCCESS}
}

func (s *systemd) Remove(name string) api.Status {
	unixname := config.APP_PREFIX + name
	errs := []string{}

	cmd := exec.Command("systemctl", "stop", unixname)
	if stdout, err := cmd.CombinedOutput(); err != nil {
		errs = append(errs, string(stdout))
	}

	cmd = exec.Command("systemctl", "disable", unixname)
	if stdout, err := cmd.CombinedOutput(); err != nil {
		errs = append(errs, string(stdout))
	}

	if err := os.Remove(SYSTEMD_OUTPUT + "/" + unixname + SYSTEMD_SERVICE_SUFFIX); err != nil {
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

func (s *systemd) Status(name string) api.Status {
	unixname := config.APP_PREFIX + name
	output := SYSTEMD_OUTPUT + "/" + unixname + SYSTEMD_SERVICE_SUFFIX
	hash := config.AppConfigGet(name, "systemd", "hash", "")

	bytes, err := ioutil.ReadFile(output)
	if err != nil {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: err.Error(),
		}
	}

	if hash != util.Hash(string(bytes)) || "" == hash {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.PROFILE_BASHRC_EXPIRED,
		}
	}

	cmd := exec.Command("systemctl", "status", unixname)
	if stdout, err := cmd.CombinedOutput(); err != nil {
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

func (s *systemd) Repair(name string) api.Status {
	return s.Create(name)
}

func (s *systemd) ListRemoved() []string {
	list := linux.GetAppList()
	result := []string{}

	files, err := ioutil.ReadDir(SYSTEMD_OUTPUT)

	if err != nil {
		panic(err.Error())
	}

	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			name = strings.TrimSuffix(name, SYSTEMD_SERVICE_SUFFIX)

			if strings.HasPrefix(name, config.APP_PREFIX) && !util.Includes(list, name) {
				result = append(result, name)
			}
		}
	}

	return result
}

func getEnvironment(name string) string {
	content := ""
	for _, handler := range module.Environ {
		kvs := handler(name)
		for _, line := range kvs {
			content += "Environment=" + line + "\n"
		}
	}

	return content
}
