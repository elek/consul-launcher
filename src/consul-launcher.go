package consullauncher

import (
	"github.com/hashicorp/consul/api"
	"path"
	"os"
	"io/ioutil"
	"os/exec"
	"time"
	"strconv"
	"syscall"
	"fmt"
	"strings"
)


type Plugin struct {
	CheckActivation func(uint64) bool
	PostIteration   func([]Entry, string)
	ProcessContent  func([]byte, *api.Client) []byte
}

var plugins = []Plugin{
	executorPlugin,
	templatePlugin,
}

type Entry struct {
	RelativePath   string
	ConsulKeyValue api.KVPair
}

func ReadConsul(dest, consul_path string, command []string) {

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}

	kv := client.KV()
	supervisor := make(chan bool)

	options := api.QueryOptions{}
	for {
		pairs, meta, error := kv.List(consul_path, &options)
		if error != nil {
			panic(error)
		}
		var changedPairs []Entry
		for _, kv := range pairs {
			if kv.ModifyIndex > options.WaitIndex {
				relativePath := kv.Key[len(consul_path):]
				content := kv.Value

				for _, plugin := range plugins {
					if plugin.CheckActivation(kv.Flags) {
						content = plugin.ProcessContent(content, client)
					}
				}
				saveFile(dest, relativePath, content)
				changedPairs = append(changedPairs, Entry{relativePath, *kv})
			}
		}

		for _, plugin := range plugins {
			plugin.PostIteration(changedPairs, dest)
		}

		if (options.WaitIndex == 0) {
			go startProcess(command, supervisor)
		} else {
			supervisor <- true
		}
		options.WaitIndex = meta.LastIndex
	}
}

func kerub(supervisor chan bool, process *os.Process) {
	signal := <-supervisor
	if (signal) {
		fmt.Println("Killing process " + strconv.Itoa(process.Pid) + " ")
		process.Kill()
	}

}
func startProcess(command[] string, supervisor chan bool) {
	retry := true
	for retry {
		var cmd *exec.Cmd
		if len(command) > 1 {
			cmd = exec.Command(command[0], command[1:]...)
		} else {
			cmd = exec.Command(command[0])
		}
		os.Stdout.Sync()
		cmd.Stdout = os.Stdout
		fmt.Println("Starting process: " + strings.Join(command, " "))
		err := cmd.Start()
		go kerub(supervisor, cmd.Process)
		err = cmd.Wait()
		if (err != nil) {
			if exitError, ok := err.(*exec.ExitError); ok {
				waitStatus := exitError.Sys().(syscall.WaitStatus)
				fmt.Println([]byte(fmt.Sprintf("Exit code: %d", waitStatus.ExitStatus())))
			} else {
				fmt.Println("Other error: " + err.Error())
			}
		} else {
			retry = !cmd.ProcessState.Success()
			fmt.Println("Process has been stopped with exit code: " + strconv.Itoa(int(cmd.ProcessState.Sys().(syscall.WaitStatus))))
		}
		time.Sleep(5 * time.Second)
	}
	os.Exit(0)
}

func saveFile(directory string, relative_path string, bytes []byte) {
	dest_file := path.Join(directory, relative_path)
	dest_dir := path.Dir(dest_file)
	err := os.MkdirAll(dest_dir, 0777)
	if (err != nil) {
		panic(err)
	}
	err = ioutil.WriteFile(dest_file, bytes, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println(dest_file + " file is written")
}

