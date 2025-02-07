package main

import (
	"GO/shared"
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
	if len(os.Args) < 3 {
		fmt.Println("Pour lancer : go run client.go <image_path> <filter_type>")
		return
	}

	//On utilise os.Args pour récupérer les arguments passés au programme
	imagePath := os.Args[1]
	filterType := os.Args[2]

	fileData, err := os.ReadFile(imagePath)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier :", err)
		return
	}

	var filter int
	switch filterType {
	case "1":
		filter = 1
	case "2":
		filter = 2
	case "3":
		filter = 3
	case "4":
		filter = 4
	default:
		fmt.Println("Filtre non reconnu. Veuillez choisir un filtre entre 1 et 4.")
		fmt.Println("Rappel : 1 - Niveaux de gris ; 2 - Détection de contours ; 3 - Netteté ; 4 - Flou gaussien")
		return
	}

	//connexion au serveur
	conn, err := net.Dial("tcp", adresse_server)
	if err != nil {
		fmt.Println("Erreur lors de la connexion au serveur :", err)
		return
	}
	defer conn.Close()

	//Préparation des données image pour l'envoi après encodage
	imgData := shared.ImageData{
		Name:       filepath.Base(imagePath),
		Data:       fileData,
		FilterType: filter,
	}

	//envoi des données image encodées
	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(imgData); err != nil {
		fmt.Println("Erreur lors de l'envoi de l'image :", err)
		return
	}

	//décodage de l'image traitée par le serveur qui est reçue par la connexion
	var processedImgData shared.ImageData
	decoder := gob.NewDecoder(conn)
	if err := decoder.Decode(&processedImgData); err != nil {
		fmt.Println("Erreur lors de la réception de l'image traitée :", err)
		return
	}

	// Création d'un sous-répertoire pour ce client (évite de tous les mélanger et d'avoir des problèmes de noms de fichiers)
	clientDir := filepath.Join("client_images", fmt.Sprintf("client_%d", time.Now().UnixNano()))
	if err := os.MkdirAll(clientDir, 0755); err != nil {
		fmt.Println("Erreur lors de la création du répertoire client :", err)
		return
	}

	//Enregistrement de l'image traitée reçue dans le répertoire propre à ce client
	outputPath := filepath.Join(clientDir, "modifiee_"+processedImgData.Name)
	if err := os.WriteFile(outputPath, processedImgData.Data, 0644); err != nil {
		fmt.Println("Erreur lors de la sauvegarde de l'image traitée :", err)
		return
	}

	fmt.Println("Image traitée sauvegardée sous :", outputPath)
}
