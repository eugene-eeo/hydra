package main

import "time"
import "bufio"
import "os"
import "net"

const waitDelay = 500 * time.Millisecond
const attempts = 3

func main() {
	for i := 0; i < attempts; i++ {
		conn, err := net.Dial("tcp", "localhost:9900")
		if err == nil {
			defer conn.Close()
			r := bufio.NewScanner(conn)
			for r.Scan() {
				_, _ = os.Stdout.Write(append(r.Bytes(), '\n'))
			}
			os.Exit(0)
		}
		time.Sleep(waitDelay)
	}
	os.Exit(1)
}
