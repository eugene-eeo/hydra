package main

import "bufio"
import "bytes"
import "os"
import "os/exec"

var pactlChange = []byte("change")

func pactlEvents(events chan string) (*os.Process, error) {
	cmd := exec.Command("pactl", "subscribe")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	go func() {
		r := bufio.NewScanner(out)
		for r.Scan() {
			if bytes.Contains(r.Bytes(), pactlChange) {
				events <- "pactl"
			}
		}
	}()
	err = cmd.Start()
	return cmd.Process, err
}
