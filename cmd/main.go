package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("hello world")
	if len(os.Args) < 3 {
		fmt.Println("Usage: ccrun run <command> <args>...")
	}
	args := os.Args
	if args[1] != "run" {
		fmt.Println("Usage: ccrun run <cmd> <args>...")
	}
	runCmd(args[2], args[3:])

}
func runCmd(command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Command execution failure: %v", err)
	}
}
