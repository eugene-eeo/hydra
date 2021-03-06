package main

import "bufio"
import "os"
import "os/exec"
import "bytes"

type nmcliProc struct{}

func (_ *nmcliProc) Run(events chan string) (*os.Process, error) {
	cmd := exec.Command("nmcli", "monitor")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	go func() {
		r := bufio.NewScanner(out)
		for r.Scan() {
			b := r.Bytes()
			// don't need to check for HasSuffix(b, "disconnected") because
			// HasSuffix(b, "connected") handles that as well
			if bytes.HasSuffix(b, []byte("connected")) || bytes.HasSuffix(b, []byte("available")) {
				events <- "nmcli"
			}
		}
	}()
	err = cmd.Start()
	return cmd.Process, err
}
