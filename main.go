package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("Couldn't read message: %s", err)
	}
	fmt.Printf("Reading Data from %s\n", inputFilePath)
	fmt.Println("=====================================")

	ch := getLinesChannel(file)

	for line := range ch {
		fmt.Println("read:", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer f.Close()
		defer close(ch)
		currentLineContents := ""
		for {
			b := make([]byte, 8)
			n, err := f.Read(b)
			if err != nil {
				if currentLineContents != "" {
					ch <- currentLineContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}
			str := string(b[:n])
			parts := strings.Split(str, "\n")
			for i := range len(parts) - 1 {
				ch <- currentLineContents + parts[i]
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return ch
}
