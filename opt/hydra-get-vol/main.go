package main

import "os"
import "bufio"
import "regexp"
import "strings"
import "os/exec"

var pctg = regexp.MustCompile("[0-9]+%")

func main() {
	cmd := exec.Command("amixer", "get", "Master")
	out, _ := cmd.StdoutPipe()
	cmd.Start()
	r := bufio.NewScanner(out)
	l := ""
	for r.Scan() {
		l = r.Text()
	}
	if strings.Contains(l, "[on]") {
		os.Stdout.WriteString(pctg.FindString(l))
	} else {
		os.Stdout.WriteString("M")
	}
}
