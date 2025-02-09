# Projet GO - Instructions

## Principe général du projet

Pour ce projet en GOLang, nous avons décidé d'appliquer le principe des goroutines à un filtrage d'images par un serveur.
Un client s'y connecte et transmet les informations suivantes : l'image à traiter et du filtre à appliquer, parmi les suivants :
- 1 : Niveaux de gris
- 2 : Détection de contours 
- 3 - Netteté  
- 4 - Flou gaussien  

### Fonctionnement du filtrage par le serveur  

Le serveur filtre l'image donnée en appliquant un Kernel correspondant au filtre. Cela se fait de manière parallèle, l'image étant découpée en 4 morceaux à traiter, chacun dans une goroutine différente.
Il renvoie ensuite l'image filtrée en indiquant au client l'emplacement à laquelle il peut la trouver.

De plus, par défaut, le serveur applique également le filtre de manière séquentielle avant de le faire en parallèle. Il utilise ce calcul pour comparer les temps d'exécution des deux méthodes.  
Si vous ne voulez pas cette comparaison, remplacez dans server.go les import comme suit :
```
import (
"GO/server/filters_sans_comparaison"
)
```  

Le serveur peut traiter plusieurs requêtes de clients à la fois, en traitant chaque client dans une goroutine qui lui est propre.  

### Différentes manière de lancer des clients

Nous avons implémenté deux versions différentes de client :
- une sans IHM, plus adaptée pour tester plusieurs clients rapidement.
- une avec IHM où toutes les informations sont détaillées.
Les manières de les exécuter sont décrites ci-dessous.  

Nous avons également créé un script bash qui permet de lancer un serveur et 4 clients à la fois, appliquant sur 4 images en sortie les différents filtres.

## Prérequis 

- **Environnement d'exécution** : Le script `server_et_clients.sh` nécessite un environnement Unix pour son exécution (par exemple, Linux ou macOS). Si vous êtes sous Windows, vous devez utiliser WSL pour pouvoir l'exécuter correctement.  
  **Note** : Les fichiers `client.go` et `server.go` peuvent être exécutés directement dans un environnement Windows sans problème.

- **Formats d'image supportés** : Le script supporte les formats d'image suivants :
  - PNG
  - JPG

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

#### Sans IHM, simple et efficace
```
cd client_sans_ihm
```
puis en remplaçant les paramètres :
```
go run client.go <image_path> <filter_type>
```
Rappel : <filter_type> est un entier, qui doit être parmi les valeurs suivantes :  1 - Niveaux de gris ; 2 - Détection de contours ; 3 - Netteté ; 4 - Flou gaussien

### Lancer un script qui lance un serveur un client par filtre

Lorsque vous êtes dans un terminal dans le répertoire du projet, lancez le script avec la commande :
```
bash server_et_clients.sh <image_path>
```
avec <image_path> le chemin d'accès vers l'image sur laquelle on souhaite appliquer les filtres.  
(Rappel : Il n'y a pas de choix de filtre ici puisque tous seront appliqués.)