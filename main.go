package main

import (
	"fmt"
	"net"
	"os"
)

func check(err error) bool {
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func main() {
	var ip string
	if len(os.Args) == 2 {
		ip = os.Args[1]
	} else {
		panic("Make sure you put an IP to run the server on.")
	}
	s := newServer()
	go s.run()
	addr := ip + ":8888"
	listener, err := net.Listen("tcp", addr)
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
