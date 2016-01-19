# stock-manager
Une simple gestion de stock via des fichiers .csv 
(',' comme séparateur, uniquement des nombres entiers, pas d'accent ni de caractère spécial). 

## Directory Tree
```
stock-manager-0.1/
  stock-manager-0.1.exe   // Lancer pour regénérer (stock, à commander, brouillons)
  bievre/                                     // Nom du stock [harcodé ...]
    l-erreurs                                 // Les éventuelles erreurs
    g-stock.csv                               // L'état actuel du stock
    g-produits.csv                            // L'évolution des produits en stock
    g-à commander.csv                         
        // Ce qu'il manque dans le stock (sur base de la valeur minimum des produits, à configurer dans c-config.csv
    c-config.csv                              
        // Configuration du système (les produits, leur nom, unité de stock, unité de Commande, coefficient [combien d'unité de stock y a-t-il dans une unité de commande], stock minimum
    brouillon-2016-01-19-n°4-commande.csv     // Commande générée
    brouillon-2016-01-19-n°3-inventaire.csv   // Inventaire généré
    brouillon-2016-01-19-n°2-sortie.csv       // Bon de sortie généré
    brouillon-2016-01-19-n°1-entrée.csv       // Bon d'entrée généré
    2016-01-19-n°3-inventaire.csv             // Inventaire
    2016-01-19-n°2-sortie.csv                 // Bon de sortie
    2016-01-19-n°1-entrée                     // Bon d'entrée
```

## Golang run
```bash
cd cmd/main/
go run main.go config.go
```

## Windows build
```bash
cd cmd/main/
GOOS=windows GOARCH=386 go build -o stock-manager-0.1.exe main.go config.go
```

