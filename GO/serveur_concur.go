package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"project/filters"
)

const (
	workerPoolSize = 10 // Number of workers in the pool
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Read filter type
	filterTypeStr, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading filter type:", err)
		return
	}
	filterTypeStr = strings.TrimSpace(filterTypeStr)
	filterType := 0
	fmt.Sscanf(filterTypeStr, "%d", &filterType)

	// Read input file path
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
		conn.Write([]byte("Error processing image\n"))
		return
	}

	conn.Write([]byte("Image processed and saved to " + outputPath + "\n"))
}

func worker(id int, jobs <-chan net.Conn) {
	for conn := range jobs {
		fmt.Printf("Worker %d processing connection\n", id)
		handleConnection(conn)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	fmt.Println("Server listening on port 8080...")

	// Create a channel to distribute connections to workers
	jobQueue := make(chan net.Conn, 100)

	// Start worker goroutines
	for i := 1; i <= workerPoolSize; i++ {
		go worker(i, jobQueue)
	}

	// Accept incoming connections and send them to the job queue
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		jobQueue <- conn
	}
}
