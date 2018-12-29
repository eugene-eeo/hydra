package main

import "net"

func server(events chan string) error {
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
					_, err := c.Write(b)
					if err != nil {
						_ = c.Close()
						conns = append(conns[:i], conns[i+1:]...)
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
	return nil
}
