// Ce dossier ne contient aucun code Go : ce fichier existe uniquement pour que
// `go build ./...` et `go test ./...` ignorent frontend/.
//
// Sans lui, Go descend dans frontend/node_modules et absorbe le paquet Go livré
// par la dépendance npm `flatted` (elle embarque du Go sans go.mod). Il apparaît
// alors dans la liste des paquets, et un jour où il ne compilerait plus, il
// casserait le build du backend. Un go.mod ici en fait un module distinct, que
// le module parent ignore entièrement.
module frontend

go 1.25
