package main

import "bufio"
import "os"
import "os/exec"
import "bytes"

var nmcliConnected = []byte("connected")
var nmcliAvailable = []byte("available")

func nmcliEvents(events chan string) (*os.Process, error) {
	cmd := exec.Command("nmcli", "monitor")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	go func() {
		r := bufio.NewScanner(out)
		for r.Scan() {
			if bytes.HasSuffix(r.Bytes(), nmcliConnected) || bytes.HasSuffix(r.Bytes(), nmcliAvailable) {
				events <- "nmcli"
			}
		}
	}()
	err = cmd.Start()
	return cmd.Process, err
}
