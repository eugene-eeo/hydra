package main

import "time"
import "fmt"
import "os"

const USAGE = `usage: hydra-timer <format>`

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, USAGE)
		os.Exit(1)
	}
	format := os.Args[1]
	ticks := time.Tick(1 * time.Second)
	prev := ""
	for {
		<-ticks
		curr := time.Now().Format(format)
		if prev != curr {
			prev = curr
			fmt.Printf("date\t%s\n", curr)
		}
	}
}
