package main

import "time"
import "os"
import "net"

func main() {
	tries := 5
	delay := 50 * time.Millisecond
	b := make([]byte, 256)

	for i := 0; i < tries; i++ {
		conn, err := net.Dial("tcp", "localhost:9900")
		if err != nil {
			time.Sleep(delay)
			delay *= 2
			continue
		}
		defer conn.Close()
		for {
			n, err := conn.Read(b)
			_, _ = os.Stdout.Write(b[:n])
			if err != nil {
				break
			}
		}
		return
	}
	os.Exit(1)
}
