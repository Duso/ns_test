package main

import (
	// #include <unistd.h>
	// #include <errno.h>
	"C"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		fmt.Println("Usage ./self_exe_demo run")
	}
}

func run() {
	cmd := exec.Command("/proc/self/exe", "child")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr {
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
		GidMappingsEnableSetgroups: false,
	}

	must(cmd.Run())
}

func child() {
	cmd := exec.Command("/bin/sc_runtime")
	cmd.Args = append(cmd.Args, "--rootfs", "/home/kpopstoyanov/simple_container/9005/mnt")
	cmd.Args = append(cmd.Args, "--workdir", "/home/plugins", "--uid", "5555", "--gid", "5555")
	cmd.Args = append(cmd.Args, "pwd")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("running %v as PID %d UID %d\n", cmd.Path, os.Getpid(), os.Getuid())
	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
