package main

import "time"
import "bufio"
import "os"
import "net"

func main() {
	attempts := 5
	delay := 50 * time.Millisecond

	for i := 0; i < attempts; i++ {
		conn, err := net.Dial("tcp", "localhost:9900")
		if err != nil {
			time.Sleep(delay)
			delay *= 2
			continue
		}
		defer conn.Close()
		r := bufio.NewScanner(conn)
		for r.Scan() {
			_, _ = os.Stdout.Write(append(r.Bytes(), '\n'))
		}
		return
	}
	os.Exit(1)
}
