package main

import (
	"os"
	"github.com/elek/consul-launcher/src"
)

func main() {
	dest := "/tmp"
	path := "conf"
	var arguments []string
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "--destination" {
			i++
			dest = os.Args[i]
		} else if arg == "--path" {
			i++
			path = os.Args[i]
		} else {
			arguments = os.Args[i:]
			break;
		}
	}
	if len(arguments) == 0 {
		panic("Usage: consul-launcher [--destination dir] [--path consul_path] any_command --with-args")
	}
	consullauncher.ReadConsul(dest, path, arguments)

}