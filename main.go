package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func shellOut(cmdString string, fail bool) {
	fmt.Printf("$ %s\n", cmdString)
	cmd := exec.Command("sh", "-c", cmdString)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		fmt.Print(stdout.String())
		fmt.Print(stderr.String())
		if fail {
			os.Exit(1)
		}
	}
}

func tmuxNewSession(session string) {
	shellOut(fmt.Sprintf("tmux new -d -s %s", session), true)
}

func tmuxAttach(session string) {
	shellOut(fmt.Sprintf("tmux attach -t %s", session), true)
}

func tmuxKillSession(session string) {
	shellOut(fmt.Sprintf("tmux kill-session -t %s", session), false)
}

func tmuxNewWindow(session string, name string) {
	shellOut(fmt.Sprintf("tmux new-window -t %s -n %s", session, name), true)
}

func tmuxSendKeys(session string, window int, pane int, command string) {
	shellOut(fmt.Sprintf("tmux send-keys -t %s:%d.%d '%s' Enter", session, window, pane, command), true)
}

func main() {
	var procfile = flag.String("file", "Procfile", "procfile location")
	var session = flag.String("session", "tmuxman", "tmux session name to use")

	flag.Parse()

	content, err := ioutil.ReadFile(*procfile)

	if err != nil {
		panic(err)
	}

	tmuxKillSession(*session)
	tmuxNewSession(*session)

	for i, line := range strings.Split(string(content), "\n") {
		tokens := strings.SplitN(line, ":", 2)

		if len(tokens) != 2 {
			continue
		}

		name, command := strings.TrimSpace(tokens[0]), strings.TrimSpace(tokens[1])
		fmt.Printf("%s -> %s \n", name, command)

		if i > 0 {
			tmuxNewWindow(*session, name)
		}
		tmuxSendKeys(*session, i+1, 0, command)
	}

	tmuxAttach(*session)
}
