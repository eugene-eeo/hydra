package main

import "bufio"
import "strings"
import "os"
import "os/exec"

func pactlEvents(events chan string) (*os.Process, error) {
	cmd := exec.Command("pactl", "subscribe")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	go func() {
		r := bufio.NewScanner(out)
		for r.Scan() {
			line := r.Text()
			if strings.Contains(line, "change") {
				events <- "pactl"
			}
		}
	}()
	err = cmd.Start()
	return cmd.Process, err
}
