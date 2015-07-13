package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
)

const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
	ExitCodeTomlNotFound
	ExitCodeTomlParseError
	ExitCodeGopathNotSet
	ExitCodeErrorParseFlag
)

type CLI struct {
	outStream, errStream io.Writer
}

type Config struct {
	Repos []string
}

func (cli *CLI) Run(args []string) int {
	var version, verbose bool

	flags := flag.NewFlagSet("goal", flag.ContinueOnError)
	flags.SetOutput(cli.errStream)
	flags.Usage = func() {
		fmt.Fprint(cli.errStream, helpText)
	}

	flags.BoolVar(&version, "version", false, "")
	flags.BoolVar(&verbose, "verbose", false, "")
	flags.BoolVar(&verbose, "v", false, "")

	// Parse all the flags
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeErrorParseFlag
	}

	// Version
	if version {
		fmt.Fprintf(cli.errStream, "%s v%s\n", Name, Version)
		return ExitCodeOK
	}

	if os.Getenv("GOPATH") == "" {
		fmt.Fprintln(cli.errStream, ColoredError("cannot download, $GOPATH not set"))
		return ExitCodeGopathNotSet
	}

	var tomlFile = "config.toml"
	if flags.NArg() > 0 {
		tomlFile = flags.Args()[0]
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

	failed := 0
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
				failed = failed + 1
			case <-doneCh:
				statusCh <- errOccurred
				break
			}
		}
	}()

	// return unix-like status code
	errOccurred := <-statusCh
	if verbose {
		all := len(conf.Repos)
		successed := all - failed
		percent := float64(successed) / float64(all) * 100
		fmt.Printf(
			"repos: %d\nsuccessed: %d, failed: %d, ok: %.1f%%\n",
			all,
			successed,
			failed,
			percent,
		)
		if int(percent) > 60 {
			return ExitCodeOK
		}
	}

	if errOccurred {
		return ExitCodeError
	}

	return ExitCodeOK
}

var helpText = `Usage: goal [options] [path]
gch is a tool to run "git status" in every $GOPATHs recursively.

Options:
--verbose, -v     View only directory path in $GOPATHs
                  without running git status.
--version         Print the version of this application
`
