package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		currentLineContents := ""
		for {
			buffer := make([]byte, 8, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					ch <- currentLineContents
					currentLineContents = ""
				}
				if errors.Is(err, io.EOF) {
					close(ch)
					f.Close()
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				close(ch)
				f.Close()
				break
			}
			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				ch <- currentLineContents + parts[i]
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()

	return ch
}
