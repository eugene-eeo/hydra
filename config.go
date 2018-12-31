package main

import "bufio"
import "encoding/json"
import "fmt"
import "io"
import "os"
import "os/exec"
import "regexp"

type Matcher struct {
	name  string
	regex *regexp.Regexp
}

type Proc struct {
	cmd      string
	args     []string
	matchers []Matcher
}

type Config struct {
	EnableNmcli bool
	EnablePactl bool
	Procs       []*Proc
}

type ProcConfig struct {
	Proc     []string    `json:"proc"`
	Matchers [][2]string `json:"match"`
}

type JSONConfig struct {
	EnableNmcli bool          `json:"nmcli"`
	EnablePactl bool          `json:"pactl"`
	ProcConfigs []*ProcConfig `json:"procs"`
}

func parseConfig(r io.Reader) (*Config, error) {
	d := json.NewDecoder(r)
	c := &JSONConfig{}
	err := d.Decode(c)
	if err != nil {
		return nil, err
	}
	cc := &Config{}
	cc.EnableNmcli = c.EnableNmcli
	cc.EnablePactl = c.EnablePactl
	cc.Procs = make([]*Proc, len(c.ProcConfigs))
	for i, pc := range c.ProcConfigs {
		if len(pc.Proc) == 0 {
			return nil, fmt.Errorf("parseConfig: proc is empty")
		}
		matchers := make([]Matcher, len(pc.Matchers))
		for j, m := range pc.Matchers {
			matchers[j] = Matcher{
				name:  m[0],
				regex: regexp.MustCompile(m[1]),
			}
		}
		cc.Procs[i] = &Proc{
			cmd:      pc.Proc[0],
			args:     pc.Proc[1:],
			matchers: matchers,
		}
	}
	return cc, nil
}

func (p *Proc) Run(events chan string) (*os.Process, error) {
	cmd := exec.Command(p.cmd, p.args...)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	go func() {
		// Expect that the process will close stdout when a signal is
		// sent to kill it.
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
	err = cmd.Start()
	return cmd.Process, err
}
