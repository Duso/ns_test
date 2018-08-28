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
	// Set upo the root fs and swap dir this will remove the mounts etc
	rootfs := "/home/kpopstoyanov/simple_container/9005/mnt"


	err := syscall.Chroot(rootfs)
	if err != nil {
		fmt.Println("Chroot error ", rootfs)
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Println("Chroot mounted ")

	// Mount proc
	err = syscall.Mount("proc", "/proc", "proc", 0, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Println("proc mounted ")

	chdir := "/home/plugins"

	err = syscall.Chdir(chdir)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Println("chdir done mounted ")


	cmd := exec.Command("/bin/sh", "-l")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(0),
			Gid: uint32(0),
		},
		GidMappingsEnableSetgroups: true,
	}

	fmt.Printf("running %v as PID %d UID %d\n", cmd.Path, os.Getpid(), os.Getuid())
	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
