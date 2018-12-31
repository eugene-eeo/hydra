package main

import "net"

func addNewline(s string) []byte {
	n := len(s)
	b := make([]byte, n+1)
	b[n] = '\n'
	copy(b, s)
	return b
}

func server(events chan string) error {
	subs := make(chan net.Conn)
	go func() {
		conns := []net.Conn{}
		for {
			select {
			case sub := <-subs:
				conns = append(conns, sub)
			case evt := <-events:
				b := addNewline(evt)
				for i := len(conns) - 1; i >= 0; i-- {
					c := conns[i]
					_, err := c.Write(b)
					if err != nil {
						_ = c.Close()
						copy(conns[i:], conns[i+1:])
						conns[len(conns)-1] = nil
						conns = conns[:len(conns)-1]
					}
				}
			}
		}
	}()
	ln, err := net.Listen("tcp", "localhost:9900")
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err == nil {
			subs <- conn
		}
	}
}
