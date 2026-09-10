package main

import (
	"strings"
	"testing"
	"time"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
)

// Le semis parle a la base, mais le plan, lui, est une valeur : ce sont ses
// proprietes qu'on fige ici. Aucun de ces tests n'ouvre de connexion, donc
// test.sh continue de tourner sans MySQL.

func TestGuardAcceptsLocalHostsOnly(t *testing.T) {
	for _, host := range []string{"", "localhost", "127.0.0.1", "MySQL", " localhost "} {
		if !IsLocalDBHost(host) {
			t.Errorf("%q devrait etre reconnu comme local", host)
		}
	}
	for _, host := range []string{"db.production.example.com", "10.0.0.4", "task-manager.rds.amazonaws.com"} {
		if IsLocalDBHost(host) {
			t.Errorf("%q ne devrait pas passer le garde-fou", host)
		}
	}
}

func TestDueAtIsAnOffsetFromTheSeedInstant(t *testing.T) {
	now := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)

	if at := DueAt(now, nil); at != nil {
		t.Errorf("une tache sans echeance devrait rester sans date, obtenu %v", at)
	}

	at := DueAt(now, in(-3*day))
	if at == nil {
		t.Fatal("une tache avec un ecart devrait porter une date")
	}
	if want := now.Add(-3 * day); !at.Equal(want) {
		t.Errorf("date %v, attendue %v", *at, want)
	}
}

// Les cles du plan sont resolues en identifiants par une map : une cle inconnue
// y vaudrait 0, et le semis echouerait sur une cle etrangere, loin de la faute.
func TestEveryKeyInThePlanNamesASeededUser(t *testing.T) {
	known := map[string]bool{}
	for _, u := range SeedUsers {
		known[u.Key] = true
	}

	check := func(key, where string) {
		if !known[key] {
			t.Errorf("%s reference la cle inconnue %q", where, key)
		}
	}

	for _, p := range SeedProjects {
		check(p.OwnerKey, "le projet "+p.Name)
		for _, key := range p.MemberKeys {
			check(key, "les membres de "+p.Name)
		}
		for _, task := range p.Tasks {
			check(task.CreatorKey, "le createur de "+task.Title)
			for _, key := range task.AssigneeKeys {
				check(key, "les assignes de "+task.Title)
			}
		}
	}
}

func TestSeededAccountsStayOnTheReservedDomain(t *testing.T) {
	// Le domaine est aussi l'etiquette du nettoyage : une adresse hors domaine
	// survivrait au passage suivant et bloquerait la reinsertion.
	for _, u := range SeedUsers {
		if !strings.HasSuffix(u.Email, "@"+SeedDomain) {
			t.Errorf("%s est hors du domaine seme, il echapperait au nettoyage", u.Email)
		}
	}
}

// Le jeu de donnees n'a d'interet que s'il rend visible chaque terme de
// CONTEXT.md. Ce test est la pour que la couverture ne se perde pas au fil des
// retouches du plan.
func TestPlanCoversTheWholeVocabulary(t *testing.T) {
	statuses := map[string]int{}
	priorities := map[string]int{}
	var overdue, dueSoon, farOff, orphan, multiAssignee, donePastDue, unassignedNoDue, emptyProjects int

	for _, p := range SeedProjects {
		if len(p.Tasks) == 0 {
			emptyProjects++
		}
		for _, task := range p.Tasks {
			statuses[task.Status]++
			priorities[task.Priority]++

			assigned := len(task.AssigneeKeys) > 0
			if len(task.AssigneeKeys) > 1 {
				multiAssignee++
			}
			if task.DueIn == nil {
				if !assigned {
					unassignedNoDue++
				}
				continue
			}

			due := *task.DueIn
			done := task.Status == models.TaskStatusDone
			switch {
			case done && due < 0:
				donePastDue++
			case done:
				// une tache terminee dont la date est a venir : sans interet ici
			case due < 0:
				overdue++
			case due <= 24*hour:
				dueSoon++
			default:
				farOff++
			}
			// Un Reminder sans assigne, c'est un Orphan.
			if !done && due <= 24*hour && !assigned {
				orphan++
			}
		}
	}

	for _, status := range models.AllowedTaskStatuses {
		if statuses[status] == 0 {
			t.Errorf("aucune tache au statut %s", status)
		}
	}
	for _, priority := range models.AllowedTaskPriorities {
		if priorities[priority] == 0 {
			t.Errorf("aucune tache de priorite %s", priority)
		}
	}

	for _, c := range []struct {
		got  int
		what string
	}{
		{overdue, "une tache Overdue"},
		{dueSoon, "une tache Due soon"},
		{farOff, "une tache dont l'echeance est lointaine"},
		{orphan, "un Orphan (Reminder sans assigne)"},
		{multiAssignee, "une tache a plusieurs assignes"},
		{donePastDue, "une tache DONE dont la date est passee (qui ne doit pas ressortir en Overdue)"},
		{unassignedNoDue, "une tache non assignee et sans echeance (non assignee, mais pas Orphan)"},
		{emptyProjects, "un projet sans aucune tache"},
	} {
		if c.got == 0 {
			t.Errorf("le plan ne contient pas %s", c.what)
		}
	}
}
