# stock-manager

[![Build Status](https://drone.io/github.com/olivier5741/stock-manager/status.png)](https://drone.io/github.com/olivier5741/stock-manager/latest)

[latest windows-386 release](https://drone.io/github.com/olivier5741/stock-manager/files)

Une simple gestion de stock via des fichiers .csv 
(',' comme séparateur, uniquement des nombres entiers, pas d'accent ni de caractère spécial). 
Le logiciel est disponible en [français](https://github.com/olivier5741/stock-manager/blob/master/cmd/main/fr-be.all.yaml) ainsi qu'en [anglais](https://github.com/olivier5741/stock-manager/blob/master/cmd/main/en-us.all.yaml).

## Directory Tree
```
stock-manager-0.2/
  l-erreurs                                 // Les éventuelles erreurs
  g-stock.csv                               // L'état actuel du stock
  g-produits.csv                            // L'évolution des produits en stock
  e-stock-manager-0.2.exe                   // Lancer pour regénérer (stock, produits, en attente)
  en attente-2016-01-30-n3-inventaire.csv  // Inventaire généré
  en attente-2016-01-30-n2-sortie.csv      // Bon de sortie généré
  en attente-2016-01-30-n1-entrée.csv      // Bon d'entrée généré
  2016-01-29-n1-inventaire.csv             // Inventaire
  2016-01-28-n2-sortie.csv                 // Bon de sortie
  2016-01-28-n1-entrée                     // Bon d'entrée
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

