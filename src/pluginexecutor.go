package consullauncher

import (
	"os/exec"
	"os"
	"path"
	"fmt"
	"github.com/hashicorp/consul/api"

)

var executorPlugin = Plugin{
	PostIteration:func(configFiles []Entry, dest string) {
		for _, configFile := range configFiles {
			kv := configFile.ConsulKeyValue
			if kv.Flags & 1 > 0 {
				fmt.Println("Executing: " + kv.Key)
				file := path.Join(dest, configFile.RelativePath)
				os.Chmod(file, os.FileMode(0755))
				cmd := exec.Command(file)
				cmd.Dir = dest
				cmd.Env = os.Environ()
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					fmt.Println("Error on executing command: " + err.Error())
				}
				os.Stdout.Sync()

			}
		}
	},
	ProcessContent:func(content []byte, consul *api.Client, supervisor chan bool) []byte {
		return content
	},
	CheckActivation:func(flag uint64) bool {
		return flag | 1 > 0
	},

}