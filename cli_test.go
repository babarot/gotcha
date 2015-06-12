package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestRunTomlNotFound(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := []string{}

	status := cli.Run(args)
	if status != ExitCodeTomlNotFound {
		t.Errorf("expected %d to eq %d", status, ExitCodeTomlNotFound)
	}
}

func TestRunTomlParseError(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}

	f, _ := ioutil.TempFile("", "invalidToml")
	content := []byte(`repos = [
	github.com/BurntSushi/toml,
]
`)
	f.Write(content)
	args := []string{
		f.Name(),
	}

	actual := cli.Run(args)
	expected := ExitCodeTomlParseError

	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestRunExitCodeOK(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}

	f, _ := ioutil.TempFile("", "validToml")
	content := []byte(`repos = [
	"github.com/BurntSushi/toml",
]
`)
	f.Write(content)
	args := []string{
		f.Name(),
	}

	actual := cli.Run(args)
	expected := ExitCodeOK

	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
