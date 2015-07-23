![](https://raw.githubusercontent.com/b4b4r07/screenshots/master/gotcha/logo.png)

[![Build Status](https://img.shields.io/travis/b4b4r07/gotcha.svg?style=flat-square)][travis]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GitHub release](http://img.shields.io/github/release/b4b4r07/gotcha.svg?style=flat-square)][release]

[travis]: https://travis-ci.org/b4b4r07/gotcha
[license]: https://raw.githubusercontent.com/b4b4r07/dotfiles/master/doc/LICENSE-MIT.txt
[release]: https://github.com/b4b4r07/gotcha/releases

`gotcha` is a simple tool that grabs Go packages

## Description

Once you've found a Go software title to your liking, you can easily install it with `gotcha`: just add the package name.


***DEMO:***

![](https://raw.githubusercontent.com/b4b4r07/screenshots/master/gotcha/demo.gif)

## Features

- **config.toml**

	It is possible to manage the go package list that you want to install by writing a TOML file such as the following:

	```toml
	repos = [
		"github.com/BurntSushi/toml",
		"github.com/BurntSushi/toml/cmd/tomlv",
		"github.com/b4b4r07/gch",
		"github.com/b4b4r07/go-pipe",
		"github.com/b4b4r07/gomi",
		# ...,
	]
	```

	[TOML](https://github.com/toml-lang/toml) is easier to read and easier to write than [JSON](https://json.org).

- **Install in parallel**

	Fast installation thanks to the parallel processing by goroutine.

## Usage

[`repos`](https://github.com/b4b4r07/gotcha/blob/master/example/config.toml#L1) that are described in `config.toml` will be install or update.

```console
$ gotcha --help
Usage: gotcha [options] [path]
gotcha is a simple tool that grabs Go packages

Options:
--verbose, -v     Cause gotcha to be verbose, showing items as they are installed.
--version         Print the version of this application
```

## Installation

![](https://raw.githubusercontent.com/b4b4r07/screenshots/master/gotcha/installation.png)

```console
$ curl -L git.io/gotcha | sh
```

If you want to go the Go way (install in GOPATH/bin) and just want the command:

```console
$ go get github.com/b4b4r07/gotcha
```

## Configuration

To customize gotcha settings:

```toml:
repos = [
	# Adding the repository to repos
    "github.com/BurntSushi/toml",
    "github.com/BurntSushi/toml/cmd/tomlv",
]

[emoji]

[emoji.verbose]
pass = ":ok_woman:"
fail = ":no_good:"

[emoji.download]
pass = ":arrow_right:"
fail = ":x:"

```

## License

[MIT](https://raw.githubusercontent.com/b4b4r07/dotfiles/master/doc/LICENSE-MIT.txt) © BABAROT (a.k.a. b4b4r07)
