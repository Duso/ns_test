package main

import (
	"cutil"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("/bin/sc_runtime")
	cmd.Args = append(cmd.Args, "--rootfs", "/home/kpopstoyanov/simple_container/9005/mnt")
	cmd.Args = append(cmd.Args, "--workdir", "/", "--uid", "0", "--gid", "0")
	cmd.Args = append(cmd.Args, "/bin/sh", "-l")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID: os.Getuid(),
				Size: 1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID: syscall.Getgid(),
				Size: 1,
			},
		},
	}
	process := cutil.NewRuntime(cmd)
	err := process.Start(nil)
	if err != nil {
		fmt.Println(err)
	}
}
