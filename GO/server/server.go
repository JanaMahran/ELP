package main

import (
	"GO/server/filters"
	"GO/shared"
	"context"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

const portString = ":9000"

var (
	clientCounter int        // Compteur global pour les clients
	clientMutex   sync.Mutex // Mutex pour protéger le compteur et éviter les race conditions (comme chaque client a accès au même compteur)
)

func init() {
	gob.Register(shared.ImageData{})
}

func main() {
	// On crée un contexte annulable, qui permettra d'interrompre proprement le server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// On gère les signaux d'arrêt (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan // Attente d'un signal
		fmt.Println("\nArrêt du serveur en cours...")
		cancel() // Annulation du contexte
	}()

	//Démarrage du serveur
	ln, err := net.Listen("tcp", portString)
	if err != nil {
		fmt.Println("Erreur au démarrage du serveur :", err)
		return
	}
	defer ln.Close()
	fmt.Printf("Le serveur écoute sur %s...\n", portString)

	var wg sync.WaitGroup

	// On accepte les connexions dans une goroutine séparée
	go func() {
		for {
			select {
			case <-ctx.Done(): // Vérification de si le contexte est annulé
				fmt.Println("Arrêt de l'acceptation de nouvelles connexions.")
				return
			default:
				conn, err := ln.Accept()
				if err != nil {
					fmt.Println("Erreur lors de l'acceptation de la connexion :", err)
					continue
				}

				wg.Add(1)
				go func() {
					defer wg.Done()
					gererClient(conn)
				}()
			}
		}
	}()

	// On attend que le contexte soit annulé
	<-ctx.Done()

	// On attend que toutes les goroutines clientes se terminent
	fmt.Println("Attente de la fin des goroutines clientes...")
	wg.Wait()
	fmt.Println("Serveur arrêté.")
}

// fonction qui traite la demande d'un client
func gererClient(conn net.Conn) {
	defer conn.Close()

	// On attribue un identifiant unique au client
	clientMutex.Lock()
	clientID := clientCounter
	clientCounter++
	clientMutex.Unlock()

	fmt.Printf("Nouveau client connecté : Client %d\n", clientID)

	// Création d'un répertoire temporaire pour ce client
	clientDir, err := os.MkdirTemp("", fmt.Sprintf("client_%d_", clientID))
	if err != nil {
		fmt.Printf("Erreur lors de la création du répertoire temporaire pour le Client %d : %v\n", clientID, err)
		return
	}
	defer os.RemoveAll(clientDir) //Une fois sa requête traité, le répertoire pourra être supprimé pour ne pas encombrer (quand tout sera fini)
	fmt.Printf("Répertoire temporaire créé pour le Client %d : %s\n", clientID, clientDir)

	//Décodage de l'image envoyée par le client à l'aide de gob
	var imgData shared.ImageData
	decoder := gob.NewDecoder(conn)
	if err := decoder.Decode(&imgData); err != nil {
		fmt.Printf("Erreur lors du décodage de l'image du Client %d : %v\n", clientID, err)
		return
	}
	fmt.Printf("Image reçue du Client %d : %s\n", clientID, imgData.Name)

	inputPath := filepath.Join(clientDir, "input_"+imgData.Name)
	if err := os.WriteFile(inputPath, imgData.Data, 0644); err != nil { //on écrit dans un fichier, avec les permission rw-r--r-- (644)
		fmt.Printf("Erreur lors de la sauvegarde de l'image du Client %d : %v\n", clientID, err)
		return
	}
	fmt.Printf("Image sauvegardée pour le Client %d : %s\n", clientID, inputPath)

	//On peut maintenant appliquer le filtre demandé à l'image reçue et l'enregistrer côté server dans outputPath
	outputPath := filepath.Join(clientDir, "output_"+imgData.Name)
	if err := filters.ApplyFilters(imgData.FilterType, inputPath, outputPath); err != nil {
		fmt.Printf("Erreur lors du traitement de l'image du Client %d : %v\n", clientID, err)
		return
	}
	fmt.Printf("Filtre appliqué pour le Client %d : %s\n", clientID, outputPath)

	//On lit ensuite l'image traitée
	processedData, err := os.ReadFile(outputPath)
	if err != nil {
		fmt.Printf("Erreur lors de la lecture de l'image traitée pour le Client %d : %v\n", clientID, err)
		return
	}

	processedImgData := shared.ImageData{
		Name: imgData.Name,
		Data: processedData,
	}

	//On envoie finalement l'image traitée au client en l'encodant avec gob
	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(processedImgData); err != nil {
		fmt.Printf("Erreur lors de l'envoi de l'image traitée au Client %d : %v\n", clientID, err)
		return
	}
	fmt.Printf("Image traitée envoyée au Client %d : %s\n", clientID, imgData.Name)
	fmt.Printf("Connexion du Client %d terminée.\n", clientID)
}
