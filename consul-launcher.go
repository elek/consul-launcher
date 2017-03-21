package main

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
)

func read_consul(dest, consul_path string, command []string) {

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
		for _, kv := range pairs {
			if kv.ModifyIndex > options.WaitIndex {
				relative_path := kv.Key[len(consul_path):]
				save_file(dest, relative_path, kv.Value)
			}
		}
		if (options.WaitIndex == 0) {
			go start_process(command, supervisor)
		} else {
			supervisor <- true
		}
		options.WaitIndex = meta.LastIndex
	}
}

func kerub(supervisor chan bool, process *os.Process) {
	signal := <-supervisor
	if (signal) {
		println("Killing process " + strconv.Itoa(process.Pid) + " ")
		process.Kill()
	}

}
func start_process(command[] string, supervisor chan bool) {
	retry := true
	for retry {
		var cmd *exec.Cmd
		if len(command) > 1 {
			cmd = exec.Command(command[0], command[1:]...)
		} else {
			cmd = exec.Command(command[0])
		}
		cmd.Stdout = os.Stdout
		println("Starting process")
		err := cmd.Start()
		go kerub(supervisor, cmd.Process)
		err = cmd.Wait()
		if (err != nil) {
			if exitError, ok := err.(*exec.ExitError); ok {
				waitStatus := exitError.Sys().(syscall.WaitStatus)
				println([]byte(fmt.Sprintf("Exit code: %d", waitStatus.ExitStatus())))
			} else {
				println("Other error " + err.Error())
			}
		}
		retry = !cmd.ProcessState.Success()

		println("Process has been stopped with exit code: " + strconv.Itoa(int(cmd.ProcessState.Sys().(syscall.WaitStatus))))
		time.Sleep(5 * time.Second)
	}
	os.Exit(0)
}

func save_file(directory string, relative_path string, bytes []byte) {
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
	println(dest_file + " file is written")
}

func main() {
	dest := "/tmp"
	path := "conf"
	var arguments []string
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "--desination" {
			i++
			dest = os.Args[i]
		} else if arg == "--path" {
			i++
			dest = os.Args[i]
		} else {
			arguments = os.Args[i:]
			break;
		}
	}

	read_consul(dest, path, arguments)

}