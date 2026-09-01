#!/usr/bin/env bash
#
# Lance les deux moitiés de la suite : le backend Go et le frontend Vitest.
#
#   ./test.sh            les deux
#   ./test.sh back       backend seulement
#   ./test.sh front      frontend seulement
#
# Depuis PowerShell : bash test.sh
#
# Le script ne s'arrête pas à la première erreur : il exécute tout, puis liste
# ce qui a échoué. Un backend cassé et un frontend cassé se voient en une seule
# exécution plutôt qu'en deux.
#
# Aucune de ces étapes n'a besoin de la base de données. Les 400 renvoyés par
# l'API sur un statut ou une priorité hors vocabulaire ne sont donc pas
# couverts ici : ils demandent un serveur lancé et un vrai token.

set -uo pipefail
cd "$(dirname "$0")"

target="${1:-all}"
failed=()

# Exécute une étape et retient son nom si elle échoue.
step() {
  local name="$1"
  shift

  printf '\n\033[1m── %s\033[0m\n' "$name"
  if "$@"; then
    return 0
  fi
  failed+=("$name")
  return 0
}

# gofmt ne renvoie pas un code d'erreur : il liste les fichiers mal formatés.
# Pas de sortie = tout est propre.
gofmt_check() {
  local out
  out="$(gofmt -l . 2>/dev/null | grep -v node_modules || true)"
  if [ -n "$out" ]; then
    echo "fichiers non formatés :"
    echo "$out"
    return 1
  fi
  echo "tous les fichiers Go sont formatés"
}

npm_step() {
  if [ ! -d frontend/node_modules ]; then
    echo "frontend/node_modules absent — lancer d'abord : cd frontend && npm install"
    return 1
  fi
  (cd frontend && npm "$@")
}

if [ "$target" = "all" ] || [ "$target" = "back" ]; then
  step "go build"  go build ./...
  step "gofmt"     gofmt_check
  step "go vet"    go vet ./...
  step "go test"   go test ./...
fi

if [ "$target" = "all" ] || [ "$target" = "front" ]; then
  step "vitest"    npm_step test
  step "eslint"    npm_step run lint
fi

if [ ${#failed[@]} -eq 0 ]; then
  printf '\n\033[32m✓ tout est vert\033[0m\n'
  exit 0
fi

printf '\n\033[31m✗ %d étape(s) en échec :\033[0m %s\n' "${#failed[@]}" "${failed[*]}"
exit 1
