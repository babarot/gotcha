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
	ExitCodeError int = 1 + iota
	ExitCodeTomlParseError
	ExitCodeTomlNotFound
)

type CLI struct {
	outStream, errStream io.Writer
}

type Config struct {
	Repos []string
}

func (cli *CLI) Run(args []string) int {
	var tomlFile = "config.toml"
	if len(args) > 0 {
		tomlFile = args[0]
	}

	// Validation for TOML
	if _, err := os.Stat(tomlFile); err != nil {
		msg := fmt.Sprintf("%s: no such file or directory", tomlFile)
		fmt.Fprintln(cli.errStream, ColoredError(msg))
		return ExitCodeTomlNotFound
	}

	// Parse TOML
	var conf Config
	if _, err := toml.DecodeFile(tomlFile, &conf); err != nil {
		fmt.Fprintln(cli.errStream, ColoredError(err.Error()))
		return ExitCodeTomlParseError
	}

	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)

	doneCh, outCh, errCh := Update(conf)
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

	// return unix-like status code
	errOccurred := <-statusCh
	if errOccurred {
		return ExitCodeError
	}
	return ExitCodeOK
}
