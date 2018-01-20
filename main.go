package main

import (
	"flag"
)

func main() {
	isInstall := flag.Bool("install", false, "install gohive")
	flag.Parse()

	if *isInstall {
		Install()
	} else {
		Web()
	}

}
