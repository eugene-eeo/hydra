package main

import "os"
import "fmt"
import "net"
import "time"

func formatAddr(port int) string {
	return fmt.Sprintf("localhost:%d", port)
}

func server(port int) error {
	ln, err := net.Listen("tcp", formatAddr(port))
	if err != nil {
		return fmt.Errorf("server: %e", err)
	}
	buffer := [1024]byte{}
	for {
		conn, err := ln.Accept()
		if err == nil {
			conn.SetReadDeadline(time.Now().Add(time.Millisecond * 50))
			n, err := conn.Read(buffer[:])
			if err == nil {
				if n == 1024 {
					n = 1023
				}
				buffer[n] = '\n'
				os.Stdout.Write(buffer[:n+1])
			}
			conn.Close()
		}
	}
}
