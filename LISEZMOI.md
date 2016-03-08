# Woutstock

<!-- I could give real cool names to the stock and the order manager on the woutstock team -->

Gestion de stock et de commande sur base de feuilles de calcul (spreadsheet google drive). 

## Gestion de stock

Un stock est un ensemble de produits dont le montant va évoluer dans le temps. 

Ce montant augmente lorsque des produits **entrent** dans le stock, représenté par un bon d'**entrée** (exemple : `2015-03-01-n1-entree` est la feuille de calcul représentant un bon d'entrée exécuté le 1er mars 2015). 

Le bon de sortie (exemple: `2015-03-01-n2-entree`, `n2` est nécessaire pour différencier/ordonner 2 actions se déroulant la même journée) est la feuille de calcul représentant un bon de sortie exécuté le 1er mars 2015.

Et finalement l'inventaire (exemple: `2015-03-01-n3-inventaire`) corrige le montant des produits du stock. 

En plus de ces type 3 feuilles de calcul, on a la mise à jour de produits stock (exemple: `2015-02-29-n1-mise à jour produit stock`). Elle n'influence pas le cours du stock, elle permet de mettre à jour la configuration du stock : les unités utilisé, les valeurs minimums pour signaler un manque (que le gestion de commande interprêtera comme une nécessité de recommander).


## Gestion des commandes

La gestion des commandes, sur base des informations fournies par la gestion de stock (les produits qu'il manque par rapport au stock minimum par exemple), va généré des commandes *en attente* par fournisseur.

> Une feuille de calcule en attente (exemple `en attente-2015-03-01-n1-commande-medipost`) est regénérée, sur base des informations fournies à Woutstock, tant que le préfix `en attente-` du nom de la feuille n'a pas été retiré.

## Le formulaire

Va-t-on vraiment m'expliquer ce qu'est un formulaire ? J'en utilise tous les jours : virement, déclaration d'impôts, ordonnance médicale, liste de courses, ... 

Juste une explication très concise alors :). Le formulaire est un support très prisé pour transmettre de l'information. Les gestionnaires (banquier pour le virement, vous-même pour la liste de course par exemple) vont traiter cette information (barrer un élément dans la liste) et éventuellement intéragir avec d'autres gestionnaires pour poursuivre l'éxécution de l'information (l'ordonnance s'échangera contre des médicaments auprès du pharmacien, il encodera l'achat [encore un formulaire!] qui sera ensuite transmis à son comptable par exemple).

Que l'on en soit conscient ou non, bien que ce soit excitant pour certains, le sujet est barbant ... Il est temps d'agir.

## Woutstock et le formulaire

Woutstock, c'est toi-même (en tant que gestionnaire de stock, gestionnaire de commande) et ton homologue digital : Ricky pour la gestion de stock et Micky pour la gestion des commandes. Et tous vous allez vous échanger des formulaires :).

### Ricky 

Qu'est-ce que Ricky peut faire pour toi ? Ricky, ton homologue digital pour la gestion des stock, peut traiter 3 types de formulaires :

* Le bon d'**entrée**, une réception de commande par exemple
* Le bon de **sortie**, une redistribution du stock pour consommation
* L'inventaire, ... un inventaire :)

Un quatrième type de formulaire s'ajoute aux trois autres : la mise à jour des produits du stock. Ce formulaire met à jour les informations sur les produits utiles pour le stock tel que le montant minimum nécessaire, les unités, ...

Jusqu'à présent, c'était ce que tu peux faire pour Ricky :). Ricky, de son côté, va parcourir tous les formulaires que tu lui as fournis et te diras ce qu'il y a théoriquement dans le stock (sous forme de document/formulaire évidemment), l'évolution du stock et des formulaires pré-remplis (pour les bons d'entrée,...) pour te faciliter la tâche. Il communique également avec Micky, le gestionnaire digital des commandes, pour lui signaler un manque de produits par exemple. 

### Micky

À temps partiel pour l'instant, Micky ne fait que pré-remplir les bons de commande pour chacun des fournisseurs à partir du formulaire mise à jour produit (différent de celui de la gestion de stock). De ton côté, tu dois modifier le bon de commande si nécessaire, mettre à jour les informations sur les produits (identifiant fournisseur, unité, prix par unité, fournisseur) ainsi que lier le formulaire bon de commande au template de commande du fournisseur (qui est joli :) et reprend par exemple les informations propres au fournisseur tel que son adresse, ...).   

## Woutstock concrètement

D'abord disponible sous forme d'un executable (`.exe`, `.bash`, ...) et intéragissant avec des fichiers csv au sein d'un système de fichiers "classique" : Woutstock Life ; Woustock se décline en une version plus "user-friendly" avec Google Spreadsheet et Google Drive : Woutstock Air. 

Les 2 versions s'appuient toutes les deux sur un système de fichiers (respectivement celui du système d'exploitation utilisé et Google Drive). Puisque un (bon) schéma vaut mieux qu'un long discours, voici la découpe dossier/fichier (et là c'est du concret :) ) : 

> ce qui se termine par `/` représente un dossier, les autres sont des formulaires (fichiers csv pour Woutstock Life et Google Spreadsheet pour Woutstock Air)

```bash
mon stock principal-ricky/
	données dashboard/ // contient les fichiers généré par Ricky 
		evolution-produit
		stock
	2016-03-06-n2-sortie // un bon de sortie au 6 mars 2016 (2ème action de la journée)
	2016-03-06-n1-sortie
	2016-03-03-n1-inventaire // un inventaire au 3 mars
	2016-03-02-n1-sortie
	2016-03-01-n3-sortie
	2016-03-01-n2-entree // un bon d'entrée au 1er mars
	2016-03-01-n1-mise à jour produit stock // une mise à jour des produits en stock au 1er mars (première action de la journée)
mes commandes-miky/
	g-bon de commande pharmacie
	g-bon de commande médipost
	g-bon de commande autre
	en attente-2015-06-01-n3-commande-pharmacie
	en attente-2015-06-01-n2-commande-médipost
	en attente-2015-06-01-n1-commande-autre
	2016-03-01-n1-mise à jour produit commande
dashboard // récapitulatif des informations importantes concernant le stock et les commandes
```

## Woutstock, rien ne se perd et tout se transforme

Une application informatique se caractérise souvent, de manière un peu simpliste, de la sorte : 

* Interface Graphique : ce que l'utilisateur voit (fenêtre, bouton, liste déroulante, et bien plus encore)
* Le core business : ce qui fait véritablement le travail (à partir des entrées, des sorties, des inventaires ; qu'y a-t-il dans mon stock)
* La persistence des données (la sauvegarde des différentes entrées, sorties, commandes, ...)

L'interface graphique est très gourmande en temps de développement (une équipe de un, ça reste une équipe de un), la technologie évolue souvent et c'est trés souvent restrictif vis-à-vis du support (android, windows laptop, web, ...). 

Étant donné que seul le core business résout véritablement un problème, Woutstock utilise un même partie tiers pour l'interface graphique et la persistence des données : système de fichier pour Woutstock Life et Google Drive pour Woutstock Air. 

> Woutstock n'est qu'un générateur : il consomme des formulaires/fichiers puis produits des formulaires/fichiers. Il **transforme**.

Le fichier/formulaire `2015-01-03-n1-entree` est à la fois une représentation visuelle d'un bon d'entrée ainsi que la donnée elle-même. Une bonne interface graphique est intuitive et n'expose à l'utilisateur que ce que l'application peut traiter. De l'autre côté, la base de données contient des données formatées et structurées

<!-- C'est lourd comme style !! -->






## Woutstock 

<!--
Parles-en à ton informaticien :)  
-->




## Nom des feuilles de calcul

Composition :

status-année-mois-jour-identifiant-type-tag1-tag2...-tagN

* Le **status** est l'état dans lequel se trouve le formulaire 

## Feuille de calcul "article"

Types concernés :

*  entree (stock)
*  sortie (stock)
* inventaire (stock)
* commande (commande)

La feuille de calcul se compose comme suit :

Produit | Valeur 1 | Unite 1 | Valeur 2 | Unite 2
--- | --- | --- | --- | ---
Abaisse-langue | 5 | unité(1) | 2 | boîte(100)
Mélolin 10x10cm | 4 | compresse(1) | 1 | boîte(100)

