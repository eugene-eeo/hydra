package main

import "time"
import "fmt"
import "os"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: hydra-timer <format>")
		os.Exit(1)
	}
	format := os.Args[1]
	prev := ""
	for {
		time.Sleep(1 * time.Second)
		curr := time.Now().Format(format)
		if prev != curr {
			prev = curr
			fmt.Printf("date\t%s\n", curr)
		}
	}
}
