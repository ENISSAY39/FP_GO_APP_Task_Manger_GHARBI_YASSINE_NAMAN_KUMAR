# Tests

Comment la suite de tests de ce dépôt est construite, comment la lancer, et
pourquoi elle s'arrête là où elle s'arrête. Ce document couvre les deux moitiés
(Go et frontend), mais insiste sur la moitié Go, celle dont le fonctionnement
est le moins visible : rien n'a été installé, rien n'a été configuré, et
pourtant `go test ./...` trouve les tests.

## 1. Ce qui a réellement été ajouté

Trois fichiers, plus un script de confort :

| Fichier | Rôle |
|---|---|
| `models/TaskModel_test.go` | fige le vocabulaire des tâches (statuts, priorités) |
| `controllers/tasksControllers_test.go` | fige les messages de refus renvoyés au client |
| `frontend/go.mod` | fait ignorer `frontend/` par la boîte à outils Go |
| `test.sh` | lance les deux moitiés de la suite en une commande |

Aucune dépendance n'a été ajoutée à `go.mod`. Aucun fichier de configuration de
test n'existe côté Go. C'est normal, et c'est le point suivant.

## 2. Pourquoi il n'y a rien à installer côté Go

En JavaScript, tester demande d'installer un lanceur (Vitest, Jest), de le
configurer, et de l'appeler par un script npm. En Go, le lanceur est dans le
langage : `go test` fait partie de la boîte à outils standard, au même titre que
`go build`. Il n'y a donc ni paquet à installer, ni fichier de configuration, ni
script à déclarer.

Ce que `go test` demande en échange, ce sont des **conventions de nommage**. Il
n'existe aucune liste de fichiers de test quelque part : la découverte est
entièrement mécanique.

### Convention 1 — le nom du fichier

Un fichier est un fichier de test si et seulement si son nom se termine par
`_test.go`.

```
models/TaskModel.go        →  code de production
models/TaskModel_test.go   →  tests de ce code
```

Conséquence directe : ces fichiers sont **exclus du binaire**. `go build ./...`
les ignore complètement ; ils ne sont compilés que par `go test`. Un test peut
donc importer des bibliothèques de test lourdes sans alourdir l'exécutable
livré.

Le fichier vit à côté du code qu'il teste, dans le même dossier — pas dans un
dossier `tests/` séparé. C'est la convention Go, et elle a un effet pratique :
voir le code et son test dans le même listing de dossier.

### Convention 2 — le nom de la fonction

Une fonction est un test si et seulement si :

- son nom commence par `Test` (majuscule ensuite : `TestFoo`, pas `Testfoo`) ;
- elle prend exactement un paramètre, de type `*testing.T` ;
- elle ne renvoie rien.

```go
func TestIsValidTaskStatus(t *testing.T) { ... }
```

`*testing.T` est le seul outil nécessaire. Il n'y a pas
d'`expect(...).toBe(...)` : on écrit un `if` ordinaire, et on signale l'échec
par une méthode de `t`.

| Méthode | Effet |
|---|---|
| `t.Errorf(...)` | marque le test en échec, **et continue** |
| `t.Fatalf(...)` | marque le test en échec, **et arrête cette fonction** |

La différence compte. `t.Errorf` dans une boucle sur des cas de test rapporte
*tous* les cas cassés en une exécution ; `t.Fatalf` s'utilise quand la suite du
test n'aurait plus de sens (par exemple si les deux messages d'erreur comparés
sont identiques, tout ce qui suit teste du vide).

### Convention 3 — le paquet

Un fichier `_test.go` déclare le **même paquet** que le code testé :

```go
// models/TaskModel_test.go
package models
```

C'est ce qui donne au test l'accès aux identifiants **non exportés** — ceux qui
commencent par une minuscule et qui sont invisibles depuis l'extérieur du
paquet. C'est exactement le cas de `controllers/tasksControllers_test.go` :
`invalidStatusMessage` et `invalidPriorityMessage` sont minuscules, donc privées
au paquet `controllers`. Être dans `package controllers` est la seule façon de
les appeler.

(Go autorise aussi `package models_test`, qui teste le paquet de l'extérieur, en
boîte noire, sans accès aux identifiants privés. Ce dépôt n'en a pas l'usage.)

### Convention 4 — le tableau de cas

Le motif dominant en Go n'est pas un test par cas, mais un tableau de cas
parcouru par une boucle (*table-driven test*) :

```go
cases := []struct {
	status string
	want   bool
}{
	{TaskStatusTodo, true},
	{"DOING", false},       // l'ancienne orthographe n'est plus acceptée
	{"in_progress", false}, // la casse compte
}

for _, tc := range cases {
	if got := IsValidTaskStatus(tc.status); got != tc.want {
		t.Errorf("IsValidTaskStatus(%q) = %v, attendu %v", tc.status, got, tc.want)
	}
}
```

Ajouter un cas = ajouter une ligne. Le message d'échec inclut l'entrée fautive,
sans quoi on saurait qu'un cas a cassé mais pas lequel.

Une variante existe avec `map[string]bool` (voir
`TestValidationBoundaryMatchesTheVocabulary`) quand les cas sont des paires
entrée/attendu sans ordre significatif. Attention : l'ordre de parcours d'une
map est aléatoire en Go — acceptable ici parce que les cas sont indépendants, à
éviter dès qu'un test dépend de la séquence.

## 3. Ce que les tests Go couvrent — et ce qu'ils ne couvrent pas

C'est la limite la plus importante à comprendre.

### `models/TaskModel_test.go`

Il fige le vocabulaire des tâches, décrit dans `CONTEXT.md`. Trois choses
distinctes y sont vérifiées :

1. **Les valeurs des constantes, pas seulement les constantes.** Le test compare
   `TaskStatusInProgress` à la chaîne littérale `"IN_PROGRESS"`. Écrire
   `if TaskStatusInProgress != TaskStatusInProgress` ne testerait rien. Ces
   valeurs sont partagées avec la base de données et avec le frontend : les
   renommer côté Go casserait silencieusement les lignes déjà stockées. Le test
   transforme ce risque silencieux en échec bruyant.
2. **La frontière accepté/refusé.** `IsValidTaskStatus` est testée sur ses cas
   positifs, mais surtout sur ses cas négatifs choisis : `"DOING"` (l'ancienne
   orthographe, abandonnée), `"in_progress"` (la casse compte),
   `" IN_PROGRESS "` (rien n'est rogné avant comparaison). Ce sont ces trois-là
   qui ont de la valeur — ils documentent des décisions.
3. **La complétude des ensembles.** `AllowedTaskStatuses` doit contenir
   exactement trois valeurs. Ajouter une constante sans l'ajouter à l'ensemble
   autorisé donnerait une valeur que le code connaît mais que l'API refuse ; le
   test attrape l'oubli.

### `controllers/tasksControllers_test.go`

Il teste les **messages de refus**, pas les handlers HTTP.

C'est un choix imposé par la structure du code : `CreateTask` et `UpdateTask`
ouvrent une connexion à la base dès leurs premières lignes. Les tester exigerait
une base de test, ou une refonte pour injecter le `*gorm.DB` — deux chantiers
hors sujet ici.

Ce qui reste testable sans base, c'est ce que le client lit réellement en cas
d'erreur. Quatre propriétés sont figées :

- le message **nomme le champ** (`status` / `priority`) ;
- il **rappelle la valeur refusée**, pour qu'on sache ce qui a été envoyé ;
- il **cite les valeurs acceptées**, en les lisant depuis
  `models.AllowedTaskStatuses` plutôt qu'en les recopiant — le test suit donc
  automatiquement tout ajout au vocabulaire ;
- les deux messages restent **distinguables** : un client envoyant statut *et*
  priorité d'un coup doit savoir lequel a été refusé.

### La zone non couverte, à dire clairement

**Le 400 renvoyé par l'API sur un statut hors vocabulaire n'est pas testé.** Les
tests vérifient que la fonction de validation dit « non » et que le message est
correct ; ils ne vérifient pas que le handler appelle bien cette fonction et
renvoie bien un `400`. Ce chaînon-là demande un serveur lancé, une base, et un
vrai token JWT.

C'est écrit en tête de `test.sh` pour que personne ne prenne un
`✓ tout est vert` pour une couverture qu'il n'est pas.

## 4. `frontend/go.mod` — le fichier le plus surprenant du lot

Un fichier `go.mod` dans un dossier qui ne contient pas une ligne de Go. Il
existe pour une raison précise.

`go build ./...` et `go test ./...` — le `...` signifie « ce dossier et tout ce
qu'il contient, récursivement » — descendaient dans `frontend/node_modules`. Or
la dépendance npm `flatted` embarque du code Go sans `go.mod`. La boîte à outils
Go l'absorbait donc comme un paquet du projet : il apparaissait dans la liste
des paquets compilés, et le jour où il cesserait de compiler, il casserait le
build du backend pour une raison sans aucun rapport avec le backend.

Poser un `go.mod` dans `frontend/` en fait un **module distinct**. Un module Go
ne descend jamais dans un sous-module : `frontend/` et tout ce qu'il contient
disparaissent d'un coup du champ de `./...`.

Le même raisonnement vaut pour `gofmt -l .`, qui ne connaît pas les modules et
parcourt les dossiers bêtement — d'où le `grep -v node_modules` dans `test.sh`.

## 5. Lancer la suite

```bash
bash test.sh          # tout
bash test.sh back     # Go seulement
bash test.sh front    # frontend seulement
```

Depuis PowerShell, il faut bien préfixer par `bash` : `./test.sh` ne
fonctionnera pas.

Le script enchaîne six étapes :

| Étape | Commande | Ce qu'elle attrape |
|---|---|---|
| `go build` | `go build ./...` | erreurs de compilation |
| `gofmt` | `gofmt -l .` | fichiers non formatés |
| `go vet` | `go vet ./...` | erreurs que le compilateur laisse passer |
| `go test` | `go test ./...` | les tests ci-dessus |
| `vitest` | `npm test` | tests du frontend |
| `eslint` | `npm run lint` | style du frontend |

Deux détails de conception du script :

- **Il ne s'arrête pas à la première erreur.** Il exécute tout, puis liste les
  étapes en échec. Un backend cassé et un frontend cassé se voient en une
  exécution plutôt qu'en deux.
- **`gofmt` demande un traitement à part.** `gofmt -l` ne renvoie pas de code
  d'erreur : il se contente de *lister* les fichiers mal formatés, et sort avec
  `0` même quand la liste n'est pas vide. Le script teste donc la sortie, pas le
  code de retour.

Les étapes Go se lancent aussi à la main, sans le script :

```bash
go test ./...                                  # tout le module
go test ./models/                              # un paquet
go test -run TestIsValidTaskStatus ./models/   # un test
go test -v ./...                               # détail test par test
```

Par défaut `go test` est silencieux quand tout passe (`ok`, une ligne par
paquet). `-v` affiche chaque test. Un paquet sans fichier de test affiche
`no test files` — ce n'est pas une erreur.

## 6. Le frontend, en un paragraphe

Vitest est installé en `devDependencies`, lancé par `npm test` (`vitest run`,
sans le mode watch), avec `jsdom` pour simuler un navigateur et
`@testing-library/react` pour interroger le rendu comme le ferait un
utilisateur. Le fichier `frontend/src/test/setup.js` charge les matchers de
`@testing-library/jest-dom`. Les tests vivent à côté des composants
(`TaskCard.test.jsx` face à `TaskCard.jsx`), comme en Go. `test.sh` refuse de
lancer cette moitié si `frontend/node_modules` est absent, et le dit plutôt que
d'échouer obscurément : `cd frontend && npm install`.

## 7. Ajouter un test Go

1. Créer (ou ouvrir) le fichier `<NomDuFichier>_test.go` à côté du code, avec le
   même `package`.
2. Écrire `func TestQuelqueChose(t *testing.T)`.
3. Préférer un tableau de cas à une suite de fonctions quasi identiques.
4. Inclure l'entrée dans le message d'échec :
   `t.Errorf("f(%q) = %v, attendu %v", in, got, want)`.
5. Choisir des cas négatifs qui documentent une décision, pas du bruit.
6. Vérifier avec `bash test.sh back`.

Si le code à tester touche la base de données, il n'est pas testable en l'état :
soit on extrait la logique pure dans une fonction séparée (comme
`invalidStatusMessage`), soit le test appartient à une couche d'intégration qui
n'existe pas encore dans ce dépôt.
