package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	const serverAddr = "localhost:42069"

	updAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		log.Fatalf("Error: net.ResolveUDPAddr: %s", err)
	}

	conn, err := net.DialUDP("udp", nil, updAddr)
	if err != nil {
		log.Fatalf("Error: net.DialUDP: %s", err)
	}
	defer conn.Close()

	fmt.Printf("Sending to %s. Type your message and press Enter to send. Press Ctrl+C to exit.\n", serverAddr)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error: reader.ReadString: %s", err)
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Printf("Error: conn.Write: %s", err)
		}

		fmt.Printf("Message sent: %s", line)
	}
}
