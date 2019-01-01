package main

import "time"
import "os"
import "net"

func main() {
	conn, err := net.DialTimeout("tcp", "localhost:9900", 500*time.Millisecond)
	if err != nil {
		os.Exit(1)
	}
	b := make([]byte, 256)
	defer conn.Close()
	for {
		n, err := conn.Read(b)
		_, _ = os.Stdout.Write(b[:n])
		if err != nil {
			return
		}
	}
}
