package main

import (
	"fmt"
	"os"

	"github.com/tsukinoko-kun/drop/internal/git"
)

func main() {
	gitWarn := make(chan string)
	p := os.Args[1]
	go git.FindUncomittedGit(p, gitWarn)
	for dir := range gitWarn {
		fmt.Println(dir)
		os.Exit(1)
	}
}
