#!/bin/bash

# Vérification que l'image nécessaire en paramètre est fournie et existe
IMAGE=$1
if [ -z "$IMAGE" ] || [ ! -f "$IMAGE" ]; then
    echo "Utilisation: $0 <image> (veuillez fournir une image existante)"
    exit 1
fi

#Démarrage du serveur en arrière plan
go run server/server.go &
SERVER_PID=$!

#On attend que le serveur soit prêt
sleep 2

#On lance ensuite plusieurs clients
go run client_sans_ihm/client.go "$IMAGE" 1 &
CLIENT1_PID=$!
go run client_sans_ihm/client.go "$IMAGE" 2 &
CLIENT2_PID=$!
go run client_sans_ihm/client.go "$IMAGE" 3 &
CLIENT3_PID=$!
go run client_sans_ihm/client.go "$IMAGE" 4 &
CLIENT4_PID=$!

#On attend que les clients terminent
wait $CLIENT1_PID $CLIENT2_PID $CLIENT3_PID $CLIENT4_PID

## On affiche un message pour prévenir de la fin des clients et demander d'arrêter le serveur manuellement
echo "Tous les clients ont terminé."
echo "Appuyez sur Ctrl+C pour arrêter le serveur"

wait $SERVER_PID
