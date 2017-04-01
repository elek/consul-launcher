package consullauncher

import (
	"os/exec"
	"os"
	"path"
	"fmt"
)

func PostIteration(configFiles []Entry, dest string) {
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
}

var executorPlugin = Plugin{
	PostIteration: PostIteration,
}