package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	// Start listening for TCP connections
	ln, err := net.Listen("tcp", portString) //":8000"
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Server is listening on port 8000...")

	for {
		// Accept a client connection
		conn, errconn := ln.Accept()
		if errconn != nil {
			fmt.Println("Error accepting connection:", errconn)
			continue
		}
		fmt.Println("Client connected:", conn.RemoteAddr())

		// Handle the client in a separate goroutine
		go gererClient(conn)
	}
}

func gererClient(conn net.Conn) {
	defer conn.Close()

	// Envoi d'un message au client pour qu'il envoie son image : A COMPLETER

	//Client envoie image

	// Réception de l'image envoyée par le client
	fmt.Println("Réception de l'image depuis le client...")
	receivedFile, err := os.Create("received_image.jpg") // Fichier temporaire pour stocker l'image reçue
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier pour l'image reçue :", err)
		return
	}
	defer receivedFile.Close()

	_, err = io.Copy(receivedFile, conn) // Copie des données reçues dans le fichier
	if err != nil {
		fmt.Println("Erreur lors de la réception de l'image :", err)
		return
	}
	fmt.Println("Image reçue avec succès.")

	//Envoi d'un message au client pour qu'il choisisse quel filtre il veut appliquer
	//ou "Lyon c'est la plus belle ville donc pas le choix, tous les filtres seront appliqués sur une image de Lyon !"

	// Read client message (e.g., file path or image data)
	reader := bufio.NewReader(conn)
	message, _ := reader.ReadString('\n')
	fmt.Println("Message from client:", message)

	// Placeholder for processing and sending a response
	fmt.Fprintln(conn, "Message received. Processing...")
	// utile ???

	// Traitement de l'image (vous pouvez appeler votre propre fonction de traitement ici)
	fmt.Println("Traitement de l'image...")
	err = processImage("received_image.jpg", "processed_image.jpg") // Traitement et sauvegarde de l'image modifiée
	if err != nil {
		fmt.Println("Erreur lors du traitement de l'image :", err)
		return
	}
	fmt.Println("Image traitée avec succès.")

	// Envoi de l'image modifiée au client
	fmt.Println("Envoi de l'image modifiée au client...")
	processedFile, err := os.Open("processed_image.jpg") // Ouvre l'image traitée
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture de l'image traitée :", err)
		return
	}
	defer processedFile.Close()

	_, err = io.Copy(conn, processedFile) // Envoi de l'image modifiée au client
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de l'image modifiée :", err)
		return
	}
	fmt.Println("Image modifiée envoyée avec succès.")
}

/* A REMPLACER / COMPLETER PAR LE PROGRAMME AVEC FILTRES, KERNEL ETC

func processImage(inputPath, outputPath string) error {
	// Ouvrir l'image source
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture de l'image source : %v", err)
	}
	defer inputFile.Close()

	// Décoder l'image
	img, _, err := image.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("erreur lors du décodage de l'image : %v", err)
	}

	// Appliquer des traitements (dans cet exemple, on sauvegarde simplement l'image sans modification)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("erreur lors de la création du fichier de sortie : %v", err)
	}
	defer outputFile.Close()

	// Réencoder l'image (format JPEG)
	err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: 90})
	if err != nil {
		return fmt.Errorf("erreur lors de l'encodage de l'image : %v", err)
	}

	return nil
}
*/
