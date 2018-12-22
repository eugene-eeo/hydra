package main

import "bufio"
import "os/exec"
import "strings"

func nmcliEvents() (chan bool, error) {
	cmd := exec.Command("nmcli", "monitor")
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
			// disconnected has the same suffix
			if strings.HasSuffix(r.Text(), "connected") {
				events <- true
			}
		}
		_ = cmd.Wait()
	}()
	return events, nil
}
