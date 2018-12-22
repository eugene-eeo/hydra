package main

import "bufio"
import "strings"
import "os/exec"

func pactlEvents() (chan bool, error) {
	cmd := exec.Command("pactl", "subscribe")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	events := make(chan bool)
	go func() {
		r := bufio.NewScanner(out)
		for r.Scan() {
			line := r.Text()
			if strings.Contains(line, "change") && strings.Contains(line, "sink") {
				events <- true
			}
		}
		_ = cmd.Wait()
	}()
	return events, nil
}
