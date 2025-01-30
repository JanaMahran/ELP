package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"project/filters"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// lire type du filtre
	filterTypeStr, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading filter type:", err)
		return
	}
	filterTypeStr = strings.TrimSpace(filterTypeStr)
	filterType := 0
	fmt.Sscanf(filterTypeStr, "%d", &filterType)

	// lire file path
	inputPath, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading input path:", err)
		return
	}
	inputPath = strings.TrimSpace(inputPath)

	// Read output file path
	outputPath, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading output path:", err)
		return
	}
	outputPath = strings.TrimSpace(outputPath)

	fmt.Printf("Received request: Filter Type=%d, Input=%s, Output=%s\n", filterType, inputPath, outputPath)

	// Apply filter
	err = func() error {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic:", r)
			}
		}()
		filters.ApplyFilters(filterType, inputPath, outputPath)
		return nil
	}()
	if err != nil {
		log.Println("Error applying filter:", err)
		return
	}

	conn.Write([]byte("Image processed and saved to " + outputPath + "\n"))
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	fmt.Println("Server listening on port 8080...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
