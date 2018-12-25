package main

import "bufio"
import "strings"
import "os/exec"

func pactlEvents(events chan string) error {
	cmd := exec.Command("pactl", "subscribe")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	go func() {
		defer cmd.Process.Kill()
		r := bufio.NewScanner(out)
		for r.Scan() {
			line := r.Text()
			if strings.Contains(line, "change") && strings.Contains(line, "sink") {
				events <- "pactl"
			}
		}
	}()
	return cmd.Start()
}
