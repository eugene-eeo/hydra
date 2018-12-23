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
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() {
		r := bufio.NewScanner(out)
		for r.Scan() {
			// disconnected has the same suffix
			if strings.HasSuffix(r.Text(), "connected") {
				events <- "nmcli"
			}
		}
		_ = cmd.Wait()
	}()
	return nil
}
