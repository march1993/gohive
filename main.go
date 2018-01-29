package main

import (
	"flag"
	"github.com/march1993/gohive/admin"
)

func main() {
	isInstall := flag.Bool("install", false, "Install gohive.")
	clearToken := flag.Bool("clear-token", false, "Clear administration token. Run installation again to generate a new token.")
	restartService := flag.Bool("restart-service", false, "Restart the service based on $SUDO_USER.")

	flag.Parse()

	if *isInstall {
		Install()
	} else if *clearToken {
		admin.SetToken("")
	} else if *restartService {
		tryRestartService()
	} else {
		Web()
	}

}
