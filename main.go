package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("error with listening TCP: %v", err)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("connection accepted")
		}

		ch := getLinesChannel(conn)

		for line := range ch {
			fmt.Printf("%s\n", line)
		}
		fmt.Println("connection has been closed")
	}

}
