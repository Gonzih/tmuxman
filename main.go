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

func tmuxSplitWindow(session string, window int, pane int, vertically bool) {
	splitKey := "h"
	if vertically {
		splitKey = "v"
	}
	shellOut(fmt.Sprintf("tmux split-window -t %s:%d.%d -%s", session, window, pane, splitKey), true)
}

func tmuxRenameWindow(session string, window int, newName string) {
	shellOut(fmt.Sprintf("tmux rename-window -t %s:%d %s", session, window, newName), true)
}

func CreateFourSplits(session string, window int) {
	tmuxSplitWindow(session, window, 0, false)
	tmuxSplitWindow(session, window, 0, true)
	tmuxSplitWindow(session, window, 1, true)
}

func main() {
	var procfile = flag.String("file", "Procfile", "procfile location")
	var session = flag.String("session", "tmuxman", "tmux session name to use")
	var useSplits = flag.Bool("splits", false, "use splits (max 4 splits per window)")

	flag.Parse()

	content, err := ioutil.ReadFile(*procfile)

	if err != nil {
		panic(err)
	}

	tmuxKillSession(*session)
	tmuxNewSession(*session)

	splitsCounter := 0
	windowsCounter := 0

	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		tokens := strings.SplitN(line, ":", 2)

		if len(tokens) != 2 {
			continue
		}

		name, command := strings.TrimSpace(tokens[0]), strings.TrimSpace(tokens[1])
		fmt.Printf("%s -> %s \n", name, command)

		if i == 0 {
			tmuxRenameWindow(*session, 1, name)
		}

		if splitsCounter == 4 {
			splitsCounter = 0
		}

		if *useSplits {
			if splitsCounter == 0 {
				if i > 0 {
					tmuxNewWindow(*session, name)
					windowsCounter = windowsCounter + 1
				}
				CreateFourSplits(*session, windowsCounter+1)
			}
			tmuxSendKeys(*session, windowsCounter+1, splitsCounter, command)
			splitsCounter = splitsCounter + 1
		} else {
			if i > 0 {
				tmuxNewWindow(*session, name)
				windowsCounter = windowsCounter + 1
			}
			tmuxSendKeys(*session, windowsCounter+1, splitsCounter, command)
		}
	}

	tmuxAttach(*session)
}
