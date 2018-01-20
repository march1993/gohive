package main

import (
	"fmt"
	"github.com/march1993/gohive/config"
)

func main() {
	fmt.Println("get_conf:" + config.Get("test", "test11"))
}
