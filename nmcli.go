package main

import "bufio"
import "os"
import "os/exec"
import "strings"

func nmcliEvents(events chan string) (*os.Process, error) {
	cmd := exec.Command("nmcli", "monitor")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	go func() {
		defer cmd.Process.Kill()
		r := bufio.NewScanner(out)
		for r.Scan() {
			// disconnected has the same suffix
			line := r.Text()
			if strings.HasSuffix(line, "connected") || strings.HasSuffix(line, "available") {
				events <- "nmcli"
			}
		}
	}()
	err = cmd.Start()
	return cmd.Process, err
}
