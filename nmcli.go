package main

import "bufio"
import "os/exec"
import "strings"

func nmcliEvents(events chan string) error {
	cmd := exec.Command("nmcli", "monitor")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	go func() {
		r := bufio.NewScanner(out)
		for r.Scan() {
			// disconnected has the same suffix
			line := r.Text()
			if strings.HasSuffix(line, "connected") || strings.HasSuffix(line, "available") {
				events <- "nmcli"
			}
		}
		_ = cmd.Wait()
	}()
	return cmd.Start()
}
