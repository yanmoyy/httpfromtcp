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
	defer file.Close()

	fmt.Printf("Reading Data from %s\n", inputFilePath)
	fmt.Println("=====================================")

	currentLineContents := ""
	for {
		b := make([]byte, 8)
		n, err := file.Read(b)
		if err != nil {
			if currentLineContents != "" {
				fmt.Printf("read: %s\n", currentLineContents)
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
			fmt.Printf("read: %s%s\n", currentLineContents, parts[i])
			currentLineContents = ""
		}
		currentLineContents += parts[len(parts)-1]
	}
}
