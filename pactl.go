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
			// ideally we want to emit events for just sink changes
			// but if we do that then the listeners will not get an
			// accurate volume read
			if bytes.Contains(r.Bytes(), []byte("change")) {
				events <- "pactl"
			}
		}
	}()
	err = cmd.Start()
	return cmd.Process, err
}
