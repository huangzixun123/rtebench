package main

import "github.com/huangzixun123/rtebench/cmd"

func main() {
	command := cmd.NewDefaultRetctlCommand()
	command.Execute()
}
