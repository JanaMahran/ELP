package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	// type du filtre
	fmt.Println("Choisis ton filtre:")
	fmt.Println("1 - Detection de bordures")
	fmt.Println("2 - Sharpen")
	fmt.Print("Entres nombre du filtre: ")
	filterType, _ := reader.ReadString('\n')
	filterType = strings.TrimSpace(filterType)

	// path image a traiter 
	fmt.Print("Enter input image path: ")
	inputPath, _ := reader.ReadString('\n')
	inputPath = strings.TrimSpace(inputPath)

	// path image traitée
	fmt.Print("Enter output image path: ")
	outputPath, _ := reader.ReadString('\n')
	outputPath = strings.TrimSpace(outputPath)

	// envoyer donnees au serveur
	fmt.Fprintln(conn, filterType)
	fmt.Fprintln(conn, inputPath)
	fmt.Fprintln(conn, outputPath)

	// lire reponse du serveur
	serverResponse, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Server response:", serverResponse)
}
