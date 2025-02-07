package main

import (
	"GO/shared"
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"
)

const adresse_server = "localhost:9000"

func init() {
	gob.Register(shared.ImageData{})
}

func main() {
	fmt.Print("Entrez le chemin du fichier image à envoyer : ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	imagePath := scanner.Text()

	fileData, err := os.ReadFile(imagePath)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier :", err)
		return
	}

	var filterType int
	for { // Boucle jusqu'à obtenir un choix valide
		fmt.Println("Entrez le numéro correspondant au filtre de votre choix parmi les suivants :")
		fmt.Println("1 - Niveaux de gris")
		fmt.Println("2 - Détection de contours")
		fmt.Println("3 - Netteté")
		fmt.Println("4 - Flou gaussien")
		fmt.Print("Votre choix : ")

		scanner.Scan()
		filterChoice := scanner.Text()

		switch filterChoice {
		case "1":
			filterType = 1
		case "2":
			filterType = 2
		case "3":
			filterType = 3
		case "4":
			filterType = 4
		default:
			fmt.Println("Choix invalide, veuillez entrer un chiffre entre 1 et 4.")
			continue // On redemande le choix sans sortir de la boucle
		}
		break // Sortie de la boucle si le choix est valide
	}

	//connexion au serveur
	conn, err := net.Dial("tcp", adresse_server)
	if err != nil {
		fmt.Println("Erreur lors de la connexion au serveur :", err)
		return
	}
	defer conn.Close() //une fois tout le reste exécuté, on fermera la connexion

	//On prépare les données image pour l'envoi après encodage
	imgData := shared.ImageData{
		Name:       filepath.Base(imagePath), // Pour utiliser uniquement le nom du fichier même si on a un chemin complet
		Data:       fileData,
		FilterType: filterType,
	}

	//envoi des données image encodées
	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(imgData); err != nil {
		fmt.Println("Erreur lors de l'envoi de l'image :", err)
		return
	}

	//On décode l'image traitée par le serveur qui est reçue par la connexion
	var processedImgData shared.ImageData
	decoder := gob.NewDecoder(conn)
	if err := decoder.Decode(&processedImgData); err != nil {
		fmt.Println("Erreur lors de la réception de l'image traitée :", err)
		return
	}

	// Création d'un sous-répertoire pour ce client (évite de tous les mélanger et d'avoir des problèmes de noms de fichiers)
	clientDir := filepath.Join("client_images", fmt.Sprintf("client_%d", time.Now().Unix()))
	//on associe au client dans le nom du répertoire l'heure exacte où il est traité
	if err := os.MkdirAll(clientDir, 0755); err != nil { // création du répertoire avec les permissions rwxr-xr-x (755)
		fmt.Println("Erreur lors de la création du répertoire client :", err)
		return
	}

	//On enregistrera l'image traitée reçue dans le répertoire propre à ce client
	outputPath := filepath.Join(clientDir, "modifiee_"+processedImgData.Name)
	if err := os.WriteFile(outputPath, processedImgData.Data, 0644); err != nil {
		fmt.Println("Erreur lors de la sauvegarde de l'image traitée :", err)
		return
	}

	fmt.Println("Image traitée sauvegardée sous :", outputPath)
}
