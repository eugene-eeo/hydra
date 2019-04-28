package main

import "bufio"
import "encoding/json"
import "fmt"
import "io"
import "os"
import "os/exec"
import "regexp"

type Runnable interface {
	Run(chan string) (*os.Process, error)
}

type Matcher struct {
	name  string
	regex *regexp.Regexp
}

type Proc struct {
	cmd      string
	args     []string
	matchers []Matcher
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

func parseConfig(r io.Reader) ([]Runnable, error) {
	d := json.NewDecoder(r)
	c := &JSONConfig{}
	if err := d.Decode(c); err != nil {
		return nil, err
	}
	return getRunnables(c)
}

func getRunnables(c *JSONConfig) ([]Runnable, error) {
	procs := make([]Runnable, len(c.ProcConfigs))
	for i, pc := range c.ProcConfigs {
		if len(pc.Proc) == 0 {
			return nil, fmt.Errorf("parseConfig: procs[%d]: proc is empty", i)
		}
		matchers := make([]Matcher, len(pc.Matchers))
		for j, pair := range pc.Matchers {
			name := pair[0]
			regex, err := regexp.Compile(pair[1])
			if err != nil {
				return nil, fmt.Errorf("parseConfig: procs[%d].match[%d]: error parsing regex", i, j)
			}
			matchers[j] = Matcher{name, regex}
		}
		procs[i] = &Proc{
			cmd:      pc.Proc[0],
			args:     pc.Proc[1:],
			matchers: matchers,
		}
	}
	if c.EnableNmcli {
		procs = append(procs, &nmcliProc{})
	}
	if c.EnablePactl {
		procs = append(procs, &pactlProc{})
	}
	return procs, nil
}

func (p *Proc) Run(events chan string) (*os.Process, error) {
	cmd := exec.Command(p.cmd, p.args...)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if len(p.matchers) == 0 {
		// just forward everything
		go func() {
			r := bufio.NewScanner(out)
			for r.Scan() {
				events <- r.Text()
			}
		}()
	} else {
		go func() {
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
