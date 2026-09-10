# Instructions du projet — FP_GO_APP_Task_Manager_GHARBI_YASSINE_NAMAN_KUMAR

Ces règles s'appliquent à toute intervention de Claude sur ce dépôt, quelle que
soit la branche.



## 2. Pas de co-authorship
- Ne jamais s'ajouter comme co-auteur dans les commits (pas de ligne
  `Co-Authored-By: Claude`, ni mention équivalente).
- Les commits doivent rester attribués uniquement à l'utilisateur.

## 3. Résumé de fin d'intervention
- À la fin de chaque session de travail sur ce projet, fournir un résumé clair
  de ce qui a été fait (fichiers modifiés/créés, décisions prises, points
  restants).

## 4. Wrapping de la prose Markdown
Voir le skill `wrap-prose` (`.claude/skills/wrap-prose/SKILL.md`), qui détaille
la règle et l'outillage.

La prose des fichiers `.md` du dépôt est repliée à **80 colonnes**, via le
script `scripts/wrap-markdown.mjs`. Rien ne replie tout seul : il faut lancer le
script à la main sur ce qui est écrit.

**Ce qui est vérifié.** La CI (`.github/workflows/docs.yml`) passe
`wrap-markdown.mjs --check` sur les `.md` du dépôt et échoue en nommant les
fichiers fautifs. Les skills vendorisés — `.agents/` et `.claude/skills/` — en
sont exclus : ils viennent de l'amont, les reformater ne ferait que polluer le
diff et serait défait à la mise à jour suivante. Seul `wrap-prose` y échappe,
puisqu'il documente cette règle et qu'il est maintenu ici. Elle ne corrige rien
et ne voit que les fichiers du dépôt : un corps de commit, d'issue ou de PR lui
échappe par construction, puisqu'il n'est jamais suivi par git.

**Ce qui ne l'est pas.** Le repli avant envoi reste une discipline. Un hook
`PreToolUse` a été tenté puis retiré : malgré plusieurs pistes écartées, il ne
s'est jamais déclenché une seule fois. La règle tient parce qu'elle est écrite
ici, pas parce qu'un automatisme y veille.

Replier également avant envoi :
- corps de commit ;
- tout bloc Markdown destiné à être copié-collé ailleurs.

**Ne jamais faire recopier un bloc de prose depuis le terminal.** Le texte
affiché est une ligne logique unique que le terminal se contente d'afficher
repliée : le copier-coller la restitue non repliée, ou repliée à la largeur de
la fenêtre. Écrire dans un fichier `.md`, le passer au script, puis passer le
fichier à l'outil :

```bash
git commit -F message.md                 # jamais un message dicté à recopier
gh pr create --body-file corps.md
gh issue comment 78 --body-file corps.md
```

**La largeur dépend de la destination — ce n'est pas 80 partout.**

| Destination | Format | Pourquoi |
|---|---|---|
| Fichier `.md` du dépôt | 80 colonnes | rendu en CommonMark (sauts simples ignorés) ; la source se lit en diff |
| Corps ou commentaire d'issue / PR GitHub | **non replié** | GitHub y active les sauts de ligne **durs** : du texte pré-replié s'affiche en escalier |
| Corps de message de commit | 72 colonnes | `git log` indente de 4 (72 + 4 = 76) ; sujet ≤ 50 |
| Teams, Slack, Discord | **non replié** | ces clients font leur propre reflow ; du texte pré-replié y sort en escalier |

Le piège : un fichier `.md` du dépôt et un commentaire GitHub sont tous les deux
du « Markdown sur GitHub », mais passent par deux moteurs de rendu différents.
Les fichiers ont des sauts souples, les champs de commentaire des sauts durs.
Replier les premiers, jamais les seconds.

Pour GitHub, Teams et consorts, il faut donc l'opération inverse : recoller
chaque paragraphe en une seule ligne logique (`--unwrap`). Le repli est
réversible sans perte, les deux sens sont toujours disponibles.

```bash
node scripts/wrap-markdown.mjs message.md              # 80 (défaut)
node scripts/wrap-markdown.mjs --width 72 message.md   # commit
node scripts/wrap-markdown.mjs --unwrap message.md     # Teams / Slack
node scripts/wrap-markdown.mjs --check *.md            # sortie 1 si non replié
```

Options complémentaires : `--stdout` (afficher sans réécrire le fichier) et
`--check` (sortie 1 en nommant les fichiers non conformes, utilisable en CI).

> `wrap-clipboard.ps1` n'existe pas dans ce dépôt. Passer par un fichier et
> `--stdout` / `--unwrap` à la place.

Ne jamais replier : blocs de code, tableaux, titres, front matter, définitions
de liens. Ne jamais couper une URL, un lien `[texte](url)` ni un span de code.

Le glossaire de `CONTEXT.md` garde chaque terme sur sa propre ligne : terminer
la ligne du terme, et celle qui précède un `_Avoid_:`, par deux espaces. C'est
un saut dur — le script le conserve, et CommonMark le rend vraiment, alors qu'un
simple retour à la ligne serait replié dans le paragraphe.

## Agent skills

### Issue tracker

Issues are tracked as GitHub issues in
`ENISSAY39/FP_GO_APP_GHARBI_YASSINE_NAMAN_KUMAR` (via `gh`). See
`docs/agents/issue-tracker.md`.

### Triage labels

Default canonical labels: `needs-triage`, `needs-info`, `ready-for-agent`,
`ready-for-human`, `wontfix`. See `docs/agents/triage-labels.md`.

### Tests

`bash test.sh` lance les deux moitiés de la suite (Go et frontend). Les
conventions de test Go, ce qui est couvert et ce qui ne l'est pas : voir
`docs/agents/tests.md`.

### Domain docs

Single-context layout: root `CONTEXT.md` + `docs/adr/`. See
`docs/agents/domain.md`.
