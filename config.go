package main

import "bufio"
import "encoding/json"
import "fmt"
import "io"
import "os/exec"

type MatcherConfig struct {
	Name    string      `json:"name"`
	Matcher interface{} `json:"matcher"`
	matcher Matcher
}

type ProcConfig struct {
	Proc     []string         `json:"proc"`
	Matchers []*MatcherConfig `json:"matchers"`
}

type Config struct {
	EnableNmcli bool          `json:"nmcli"`
	EnablePactl bool          `json:"pactl"`
	ProcConfigs []*ProcConfig `json:"procs"`
}

func parseConfig(r io.Reader) (*Config, error) {
	d := json.NewDecoder(r)
	c := &Config{}
	err := d.Decode(c)
	if err != nil {
		return nil, err
	}
	for _, pc := range c.ProcConfigs {
		if len(pc.Proc) == 0 {
			return nil, fmt.Errorf("parseConfig: proc is empty")
		}
		for _, mc := range pc.Matchers {
			m, err := interfaceToMatcher(mc.Matcher)
			if err != nil {
				return nil, fmt.Errorf("parseConfig: error parsing \"%s\": %e", mc.Name, err)
			}
			mc.matcher = m
		}
	}
	return c, nil
}

func (p *ProcConfig) Run(events chan string) error {
	cmd := exec.Command(p.Proc[0], p.Proc[1:]...)
	out, _ := cmd.StdoutPipe()
	go func() {
		// Expect that the process will close stdout when a signal is
		// sent to kill it.
		r := bufio.NewScanner(out)
		for r.Scan() {
			line := r.Text()
			for _, mc := range p.Matchers {
				if mc.matcher.Match(line) {
					events <- mc.Name
					break
				}
			}
		}
		_ = cmd.Wait()
	}()
	return cmd.Start()
}
