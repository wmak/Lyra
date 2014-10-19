package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("/home/terrence/develop/Lyra/testing/library_client.go", "*.go")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}