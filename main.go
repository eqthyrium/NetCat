package main

import (
	"fmt"
	"os"

	"TCPChat/cmd"
)

func main() {
	err := cmd.FromMain(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
