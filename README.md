# Tmuxman [![Build Status](https://travis-ci.org/Gonzih/tmuxman.svg?branch=master)](https://travis-ci.org/Gonzih/tmuxman)

Clone of Foreman for tmux written in golang.

## Installation

```bash
go get github.com/Gonzih/tmuxman
```

## Usage

```bash
$ tmuxman --help
Usage of tmuxman:
  -file string
        procfile location (default "Procfile")
  -session string
        tmux session name to use (default "tmuxman")
  -splits
        use splits (max 4 splits per window)

$ tmuxman --file Procfile.development --session myproject --splits
```

## License

MIT

## Authors

Max Gonzih
