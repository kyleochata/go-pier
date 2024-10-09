package main

import (
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
		panic("pass me an argument please")
	}
}

func run() {
	fmt.Printf("Running %v as pid %d\n", os.Args[2:], os.Getpid())
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS, // New UTS and PID namespace
	}
	must(cmd.Run())
}

func child() {
	fmt.Printf("Running %v\n", os.Args[2:])
	fmt.Printf("Running %v as pid %d\n", os.Args[2:], os.Getpid())
	// Change hostname within the new namespace
	if err := syscall.Sethostname([]byte("container")); err != nil {
		panic(err)
	}
	syscall.Chroot("/vagrant/ubuntu-fs")
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "") //Mount proc. Without this ps will return library fatal, lookup self: host: (wsl -d Ubuntu) --> wsl:  go run main.go run /bin/bash --> Container: ps; ps: output: library fatal, lookup self

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())
	must(syscall.Unmount("proc", 0))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
