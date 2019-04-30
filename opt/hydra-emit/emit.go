package main

import "fmt"
import "net"

func emit(port int, content []byte) error {
	conn, err := net.Dial("tcp", formatAddr(port))
	if err != nil {
		return fmt.Errorf("emit: %e", err)
	}
	_, err = conn.Write(content)
	if err != nil {
		return fmt.Errorf("emit: %e", err)
	}
	return nil
}
