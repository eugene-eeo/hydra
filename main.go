package main

import "os"
import "os/signal"
import "syscall"
import "github.com/mitchellh/go-homedir"

func read_config() *Config {
	path, err := homedir.Expand("~/.hydrarc.json")
	must(err)
	f, err := os.Open(path)
	must(err)
	defer f.Close()
	config, err := parseConfig(f)
	must(err)
	return config
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	config := read_config()
	events := make(chan string, 5)
	procs := []*os.Process{}
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for _ = range sigs {
			for _, proc := range procs {
				_ = proc.Kill()
			}
			os.Exit(0)
		}
	}()

	if config.EnablePactl {
		proc, err := pactlEvents(events)
		must(err)
		procs = append(procs, proc)
	}
	if config.EnableNmcli {
		proc, err := nmcliEvents(events)
		must(err)
		procs = append(procs, proc)
	}
	for _, p := range config.Procs {
		proc, err := p.Run(events)
		must(err)
		procs = append(procs, proc)
	}
	server(events)
}
