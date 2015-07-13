package main

import "os"

const (
	Name    = "goal"
	Version = "0.1.1"
)

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}
