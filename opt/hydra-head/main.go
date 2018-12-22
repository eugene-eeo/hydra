package main

import "bufio"
import "os"
import "net"

func main() {
	conn, err := net.Dial("tcp", "localhost:9900")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	r := bufio.NewScanner(conn)
	for r.Scan() {
		os.Stdout.WriteString(r.Text() + "\n")
	}
}
