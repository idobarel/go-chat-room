package main

import (
	"fmt"
	"net"
)

func check(err error) bool {
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func main() {
	s := newServer()
	go s.run()

	listener, err := net.Listen("tcp", "192.168.1.126:8888")
	if !check(err) {
		panic("Error Creating Server.")
	}

	defer listener.Close()
	fmt.Println("Started Serted Server On :8888.")

	for {
		conn, err := listener.Accept()
		if !check(err) {
			continue
		}
		go s.newClient(conn)
	}
}
