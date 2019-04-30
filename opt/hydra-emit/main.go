package main

import "os"
import "fmt"

const Port = 9901
const Usage = `usage:
	%s emit <event>
	%s server
`

func die_usage() {
	fmt.Fprintf(os.Stderr, Usage, os.Args[0], os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) == 1 {
		die_usage()
	}
	var err error
	switch os.Args[1] {
	case "emit":
		if len(os.Args) != 3 {
			die_usage()
		}
		err = emit(Port, []byte(os.Args[2]))
	case "server":
		err = server(Port)
	default:
		die_usage()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%e\n", err)
		os.Exit(1)
	}
}
