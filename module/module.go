package module

import (
	"github.com/march1993/gohive/api"
)

type Module interface {
	Create(name string) api.Status // create application
	Remove(name string) api.Status // remove application
	Status(name string) api.Status // show status
	Repair(name string) api.Status // Repair an application when it's broken or copied to the hive directory

	ListRemoved() []string // List application names those have been removed
}

type moduleItem struct {
	name   string
	module Module
}

var Modules = []moduleItem{}
var Runcom = []func(name string) []string{}

func RegisterRuncom(handler func(name string) []string) {
	Runcom = append(Runcom, handler)
}

func RegisterModule(name string, module Module) {

	Modules = append(Modules, moduleItem{name: name, module: module})

}

func StatusApp(app string) map[string]api.Status {
	ret := map[string]api.Status{}

	for _, item := range Modules {
		ret[item.name] = item.module.Status(app)
	}

	return ret
}

func CreateApp(app string) map[string]api.Status {
	ret := map[string]api.Status{}

	for _, item := range Modules {
		ret[item.name] = item.module.Create(app)
	}

	return ret
}

func RemoveApp(app string) map[string]api.Status {
	ret := map[string]api.Status{}

	for _, item := range Modules {
		ret[item.name] = item.module.Remove(app)
	}

	return ret
}

func RepairApp(app string) map[string]api.Status {
	ret := map[string]api.Status{}

	for _, item := range Modules {
		ret[item.name] = item.module.Repair(app)
	}

	return ret
}

func ListRemovedApp() map[string][]string {
	ret := map[string][]string{}

	for _, item := range Modules {
		ret[item.name] = item.module.ListRemoved()
	}

	return ret
}
