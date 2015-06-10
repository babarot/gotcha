package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
)

const (
	ExitCodeOK    int = 0
	ExitCodeError     = 1 + iota
	ExitCodeTomlParseError
	ExitCodeTomlNotFound
	ExitCodeBadArgs
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var tomlFile = "config.toml"
	if len(args) > 0 {
		tomlFile = args[0]
	}

	if _, err := os.Stat(tomlFile); err != nil {
		msg := fmt.Sprintf("%s: no such file or directory", tomlFile)
		fmt.Fprintln(cli.errStream, ColoredError(msg))
		return ExitCodeTomlNotFound
	}

	var conf Config
	if _, err := toml.DecodeFile(tomlFile, &conf); err != nil {
		fmt.Fprintln(cli.errStream, ColoredError(err.Error()))
		return ExitCodeTomlParseError
	}

	// Use CPU efficiently
	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)

	// Start releasing
	doneCh, outCh, errCh := Update(conf)

	// Receive messages
	statusCh := make(chan bool)
	go func() {
		errOccurred := false
		for {
			select {
			case out := <-outCh:
				fmt.Fprintf(cli.outStream, out)
			case err := <-errCh:
				fmt.Fprintf(cli.errStream, ColoredError(err))
				errOccurred = true
			case <-doneCh:
				statusCh <- errOccurred
				break
			}
		}
	}()

	// If more than one error is occured, return non-zero value
	errOccurred := <-statusCh
	if errOccurred {
		return ExitCodeError
	}

	return ExitCodeOK
}
