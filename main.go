package main

import (
	"fmt"
	"io"
	"os"
)

func main() {

	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	defer file.Close()

	buffer := make([]byte, 8)
	pickindex := 0

	for {
		bytesread, err := file.Read(buffer)

		if bytesread > 0 {
			chunk := buffer[:bytesread]
			fmt.Printf("read: %s\n", string(chunk))
			pickindex++
		}

		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Printf("error: %v", err)
		}
	}

}
