package main

import (
	"fmt"
	"log"
	"os"
)

const inputFilePath = "messages.txt"

func main() {
	f, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("could not open %s: %s\n", inputFilePath, err)
	}

	fmt.Printf("Reading data from %s\n", inputFilePath)
	fmt.Println("=====================================")

	ch := getLinesChannel(f)

	for line := range ch {
		fmt.Printf("read: %s\n", line)
	}
}
