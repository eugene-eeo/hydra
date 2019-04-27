package main

import "bufio"
import "bytes"
import "os"
import "os/exec"

type pactlProc struct{}

func (_ *pactlProc) Run(events chan string) (*os.Process, error) {
	cmd := exec.Command("pactl", "subscribe")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	go func() {
		r := bufio.NewScanner(out)
		for r.Scan() {
			if bytes.HasPrefix(r.Bytes(), []byte("Event 'change' on sink")) {
				events <- "pactl"
			}
		}
	}()
	err = cmd.Start()
	return cmd.Process, err
}
