package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	lAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Printf("error resolving local host: %v", err)
	}

	UDPconn, err := net.DialUDP("udp", nil, lAddr)
	if err != nil {
		fmt.Printf("error dialing UDP: %v", err)
	}

	defer UDPconn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("error reading string: %v", err)
		}
		b := []byte(line)
		_, err = UDPconn.Write(b)
		if err != nil {
			fmt.Printf("error writing to udp: %v", err)
		}
	}
}
