package main

import "github.com/google/shlex"
import "fmt"
import "encoding/json"
import "io"
import "regexp"

type ProcConfig struct {
	Proc     string      `json:"proc"`
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
		cmd, err := shlex.Split(pc.Proc)
		if err != nil {
			return nil, fmt.Errorf("procs[%d]: cannot parse proc: %e", i, err)
		}
		if len(cmd) == 0 {
			return nil, fmt.Errorf("procs[%d]: proc is empty", i)
		}
		matchers := make([]Matcher, len(pc.Matchers))
		for j, pair := range pc.Matchers {
			name := pair[0]
			regex, err := regexp.Compile(pair[1])
			if err != nil {
				return nil, fmt.Errorf("procs[%d].match[%d]: error parsing regex: %e", i, j, err)
			}
			matchers[j] = Matcher{name, regex}
		}
		procs[i] = &Proc{cmd, matchers}
	}
	if c.EnableNmcli {
		procs = append(procs, &nmcliProc{})
	}
	if c.EnablePactl {
		procs = append(procs, &pactlProc{})
	}
	return procs, nil
}
