package main

import "bufio"
import "strings"
import "os"
import "os/exec"

func pactlEvents(events chan string, procs *[]*os.Process) error {
	cmd := exec.Command("pactl", "subscribe")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	go func() {
		r := bufio.NewScanner(out)
		for r.Scan() {
			line := r.Text()
			if strings.Contains(line, "change") && strings.Contains(line, "sink") {
				events <- "pactl"
			}
		}
	}()
	err = cmd.Start()
	*procs = append(*procs, cmd.Process)
	return err
}
