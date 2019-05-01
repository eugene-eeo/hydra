package main

import "os"
import "os/exec"
import "bufio"
import "regexp"

type Runnable interface {
	Run(chan string) (*os.Process, error)
}

type Matcher struct {
	name  string
	regex *regexp.Regexp
}

type Proc struct {
	cmd      []string
	matchers []Matcher
}

func (p *Proc) Run(events chan string) (*os.Process, error) {
	cmd := exec.Command(p.cmd[0], p.cmd[1:]...)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if len(p.matchers) == 0 {
		// just forward everything
		go func() {
			defer out.Close()
			r := bufio.NewScanner(out)
			for r.Scan() {
				events <- r.Text()
			}
		}()
	} else {
		go func() {
			defer out.Close()
			r := bufio.NewScanner(out)
			for r.Scan() {
				b := r.Bytes()
				for _, m := range p.matchers {
					if m.regex.Match(b) {
						events <- m.name
						break
					}
				}
			}
		}()
	}
	err = cmd.Start()
	return cmd.Process, err
}
