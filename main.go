package main

import "fmt"
import "os"
import "os/signal"
import "syscall"
import "github.com/mitchellh/go-homedir"

func read_config() []Runnable {
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
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func kill(procs []*os.Process) {
	for _, proc := range procs {
		if proc != nil {
			_ = proc.Kill()
		}
	}
}

func spawn_and_listen(runnables []Runnable) error {
	events := make(chan string, 5)
	procs := make([]*os.Process, len(runnables))
	sigs := make(chan os.Signal)
	signal.Notify(sigs,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGSTOP,
	)
	go func() {
		<-sigs
		kill(procs)
		os.Exit(0)
	}()
	defer kill(procs)
	for i, p := range runnables {
		proc, err := p.Run(events)
		if err != nil {
			return err
		}
		procs[i] = proc
	}
	// gc
	runnables = nil
	return server(events)
}

func main() {
	must(spawn_and_listen(read_config()))
}
