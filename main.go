package main

import "os"
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
	if config.EnablePactl {
		must(pactlEvents(events))
	}
	if config.EnableNmcli {
		must(nmcliEvents(events))
	}
	for _, pc := range config.ProcConfigs {
		must(pc.Run(events))
	}
	server(events)
}
