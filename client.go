package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	// On demande à l'utilisateur le chemin sur sa machine de l'image à envoyer
	fmt.Print("Entrez le chemin du fichier image (jpg, png ou gif) à envoyer : ")
	//On récupère ainsi l'information entrée par l'utilisateur, conservée dans la variable imagePath :
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	imagePath := scanner.Text()

	// On vérifie ensuite si le fichier existe
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Printf("Erreur : Impossible d'ouvrir le fichier %s\n", imagePath)
		return
	}
	defer file.Close() //rappel : il y a un defer donc ne s'applique que quand le reste a été exécuté et n'a pas planté

	// Connexion au serveur TCP
	fmt.Println("Connexion au serveur...")
	adresse_server = "localhost:8000"
	conn, err := net.Dial("tcp", adresse_server) // Bien s'assurer que le serveur écoute sur ce port
	//si on a bien lancé le programme server.go avant c'est bon
	if err != nil {
		fmt.Println("Erreur lors de la connexion au serveur :", err)
		return
	}
	defer conn.Close()

	// Envoi du fichier au serveur
	fmt.Println("Envoi de l'image au serveur...")
	_, err = io.Copy(conn, file) //envoie le contenu file (flux sortant) par conn la connection TCP avec le serveur (flux entrant)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de l'image :", err)
		return
	}
	fmt.Println("Image envoyée avec succès.")

	// Réception de l'image modifiée depuis le serveur
	fmt.Println("Réception de l'image modifiée...")
	outputFile, err := os.Create("modified_image.jpg")
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier de sortie :", err)
		return
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, conn)
	if err != nil {
		fmt.Println("Erreur lors de la réception de l'image :", err)
		return
	}
	fmt.Println("Image modifiée reçue et sauvegardée sous 'modified_image.jpg'.")
}
