# stock-manager
Une simple gestion de stock via des fichiers .csv 
(';' comme séparateur, uniquement des nombres entiers, pas d'accent ni de caractère spécial). 

## Directory Tree
```
stock-manager-0.1
  input/
    2015-01-09-bievre-1-in.csv    // Bon d'entrée pour le stock de bièvre
    2015-01-09-libin-1-in.csv     // Bon d'entrée pour le stock de libin
    2015-01-09-libin-2-out.csv    // Bon de sortie pour le stock de libin (2ème action de la journée)
    2015-01-11-bievre-1-out.csv   // Bon de sortie pour le stock de bièvre
    2015-01-11-bievre-2-inv.csv   // Inventaire du stock de bièvre (2ème action de la journée)
    2015-01-12-libin-1-inv.csv    // Inventaire du stock de libin
  output/
    stock.csv             // L'état actuel du stock
    missing.csv           // Ce qu'il manque dans le stock (sur base de la valeur minimum des produits, à configurer dans *config.yaml*
  config.yaml             // Fichier de configuration du système (les produits, leur nom *name*, valeur minimum *min*, le nombre par boîte *bulk*
  log                     // Fichier de log, pour consulter les éventuelles erreurs
  stock-manager-0.1.exe   // Lancer pour regénérer les output à partir des input
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

