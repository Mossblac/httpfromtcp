package main

import (
	"fmt"
	"httpfromtcp/internal/request"
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
		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("error with parsing: %v", err)
		}
		template := `
		Request line: 
		- Method: %v
		- Target: %v
		- Version: %v
		`
		output := fmt.Sprintf(
			template,
			req.RequestLine.Method,
			req.RequestLine.RequestTarget,
			req.RequestLine.HttpVersion,
		)
		fmt.Print(output)
		fmt.Println("Headers:")
		for key, value := range req.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}

		fmt.Println("connection has been closed")
	}

}
