# Projet GO - Instructions

## Principe général du projet



## Lancer le projet

### Accéder au répertoire GO

Ouvrez un terminal et naviguez jusqu'au répertoire contenant votre projet Elm en utilisant la commande suivante :

```
cd chemin/vers/le/répertoire/GO
```

### Démarrer un serveur

Rendez-vous dans le répertoire dédié au serveur avec la commande :
```
cd server
```
Lancez le serveur qui traitera les images avec la commande suivante :
```
go run server.go
```

### Démarrer un client

Une fois un serveur lancé, on peut maintenant lancer un client qui demandera de filtrer une image. 
Commencez par ouvrir un nouveau terminal et accédez au répertoire GO.

Deux versions du client existent selon les besoins.
Appliquez l'une des deux séries de commandes suivantes (selon si vous avez du temps devant vous ou non) :

#### Pour avoir l'IHM qui demande le filtre choisi et sur quel image l'appliquer :
```
cd client_avec_ihm
```
puis
```
go run client.go
```

#### Sans IHM, simple efficace
```
cd client_avec_ihm
```
puis en remplaçant les paramètres :
```
go run client.go <image_path> <filter_type>
```
Rappel : <filter_type> est un entier, qui doit être parmi les valeurs suivantes :  1 - Niveaux de gris ; 2 - Détection de contours ; 3 - Netteté ; 4 - Flou gaussien

