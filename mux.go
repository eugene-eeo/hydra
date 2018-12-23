package main

import "os"
import "fmt"
import "net"

func server(events chan string) {
	subs := make(chan net.Conn)
	go func() {
		conns := []net.Conn{}
		for {
			select {
			case sub := <-subs:
				conns = append(conns, sub)
			case evt := <-events:
				b := append([]byte(evt), '\n')
				for i := len(conns) - 1; i >= 0; i-- {
					c := conns[i]
					if _, err := c.Write(b); err != nil {
						_ = c.Close()
						conns = append(conns[:i], conns[i+1:]...)
					}
				}
			}
		}
	}()
	ln, err := net.Listen("tcp", "localhost:9900")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for {
		if conn, err := ln.Accept(); err == nil {
			subs <- conn
		}
	}
}
