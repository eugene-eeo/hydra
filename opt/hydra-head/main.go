package main

import "io"
import "time"
import "os"
import "net"

func main() {
	tries := 5
	delay := 50 * time.Millisecond

	for i := 0; i < tries; i++ {
		conn, err := net.Dial("tcp", "localhost:9900")
		if err != nil {
			time.Sleep(delay)
			delay *= 2
			continue
		}
		defer conn.Close()
		_, _ = io.Copy(os.Stdout, conn)
		return
	}
	os.Exit(1)
}
