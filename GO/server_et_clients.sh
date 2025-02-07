#!/bin/bash

#Démarrage du serveur en arrière plan
go run server/server.go &
SERVER_PID=$!

#On attend que le serveur soit prêt
sleep 3

#On lance ensuite plusieurs clients
go run client_sans_ihm/client.go client_sans_ihm/lyon.jpg 1 &
CLIENT1_PID=$!
go run client_sans_ihm/client.go client_sans_ihm/lyon.jpg 2 &
CLIENT2_PID=$!
go run client_sans_ihm/client.go client_sans_ihm/lyon.jpg 3 &
CLIENT3_PID=$!
go run client_sans_ihm/client.go client_sans_ihm/lyon.jpg 4 &
CLIENT4_PID=$!

#On attend que les clients terminent
wait $CLIENT1_PID $CLIENT2_PID $CLIENT3_PID $CLIENT4_PID

## On affiche un message pour prévenir de la fin des clients et demander d'arrêter le serveur manuellement
echo "Tous les clients ont terminé."
echo "Appuyez sur Ctrl+C pour arrêter le serveur"

wait $SERVER_PID
