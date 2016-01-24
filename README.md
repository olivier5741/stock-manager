# stock-manager
[![Build Status](https://drone.io/github.com/olivier5741/stock-manager/status.png)](https://drone.io/github.com/olivier5741/stock-manager/latest)
Une simple gestion de stock via des fichiers .csv 
(',' comme séparateur, uniquement des nombres entiers, pas d'accent ni de caractère spécial). 
Le logiciel est disponible en [français](https://github.com/olivier5741/stock-manager/blob/master/cmd/main/fr-be.all.yaml) ainsi qu'en [anglais](https://github.com/olivier5741/stock-manager/blob/master/cmd/main/en-us.all.yaml).

## Directory Tree
```
stock-manager-0.2/
  stock-manager-0.2.exe   // Lancer pour regénérer (stock, à commander, brouillons)
  bievre/                                     // Nom du stock [harcodé ...]
    l-erreurs                                 // Les éventuelles erreurs
    g-stock.csv                               // L'état actuel du stock
    g-produits.csv                            // L'évolution des produits en stock
    g-à commander.csv                         
        // Ce qu'il manque dans le stock (sur base de la valeur minimum des produits, à configurer dans c-config.csv
    en attente-2016-01-19-n°4-commande.csv     // Commande générée
    en attente-2016-01-19-n°3-inventaire.csv   // Inventaire généré
    en attente-2016-01-19-n°2-sortie.csv       // Bon de sortie généré
    en attente-2016-01-19-n°1-entrée.csv       // Bon d'entrée généré
    c-config.csv                              
        // Configuration du système (les produits, leur nom, unité de stock, unité de Commande, coefficient [combien d'unité de stock y a-t-il dans une unité de commande], stock minimum
    2016-01-19-n°3-inventaire.csv             // Inventaire
    2016-01-19-n°2-sortie.csv                 // Bon de sortie
    2016-01-19-n°1-entrée                     // Bon d'entrée
```

## Golang run
```bash
cd cmd/main/
go run main.go
```

## Windows build
```bash
cd cmd/main/
GOOS=windows GOARCH=386 go build -o stock-manager-0.2.exe main.go
```

