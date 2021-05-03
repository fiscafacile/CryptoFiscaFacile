Pour contribuer, dans la majorité des cas vous allez vouloir ajouter le support pour une plateforme manquante. Je vous porposer de prendre comme base la Source Crypto.com existante qui est la plus aboutie.

Il vous faudra :

- [x] copier tout le répertoire, épurer ce dont vous avez pas besoin (ex: pas d'API, CSV uniquement, ou l'inverse) et renomer pour votre plateforme.

- [x] adapter le code pour votre format de CSV ou votre API.

- [x] vérifier que vos TXs n'aient aucun `Amount` négatif.

- [x] vérifier que les `Items` ne soient que des `To`/`From`/`Fee` (il faut raisoner aux bornes de la Source en question, pas du portefeuille global, faire abstraction des autres Sources).

- [x] vérifier que les `Deposits` n'aient pas de `From` (la source d'un dépot n'est pas à intégrer dans votre portefeuille relatif à la Source en cours, sauf si elle vous appartient auquel cas ce n'est pas un `Withdrawals` mais un `Transfers` qu'il faut faire).

- [x] vérifier que si vos `Deposits` ont des `Fee`, ils aient bien été payés par vous (souvent c'est la source du dépot qui paye les frais).

- [x] vérifier que les `Withrawals` n'aient pas de `To` (la destination d'un retrait n'est pas à intégrer dans votre portefeuille relatif à la Source en cours, sauf si elle vous appartient auquel cas ce n'est pas un `Deposits` mais un `Transfers` qu'il faut faire).

- [x] vérifier que les `Transfers` aient une balance nulle (frais compris).

- [x] si possible faites des tests unitaires pour vérifier votre module uniquement.

- [x] ajouter votre Source dans la doc du README.md (prendre exemple sur un autre).

- [x] ajouter le support de votre Source dans le main.go (ajouter le flag, la création de l'instance, la récupération des TXs (ParseCSV ou GetAPI) puis intégration des TXs par catégorie au portefeuille global).

- [x] faire un test d'ensemble en ne fournissant que votre source à l'outil compilé

- [x] contrôler que la balance correspond à ce que vous attendiez (attention à la date qui par défaut est au 1 Jan 2021), si une ou plusieurs balance ne correspond pas, il faut isoler un coin et suivre toutes ses transactions avec les options `-txs_display` et `-curr_filter`.

- [x] faire un `-check` pour vérifier que vos TXs répondent au critères de l'outil.

- [x] refaire un test d'ensemble avec une autre Source pour voir si des `Deposits` et des `Withdrawals` correspondants sont fusionés en `Transfers`. Si ce n'est pas le cas, vérifier les montant (doivent correspondre) ou les dates (doivent etre à 24h près).

- [ ] bien faire un go fmt, puis git commit / push et enfin demander le Pull Request.

Pour toute question, vous pouvez venir sur le groupe telegram donné en bas de la doc officielle.
