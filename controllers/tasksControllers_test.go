package controllers

import (
	"strings"
	"testing"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
)

// Les handlers eux-mêmes parlent à la base : ce qui est testable sans base,
// c'est le message de refus, et c'est lui que le client lit. On vérifie qu'il
// nomme le champ, la valeur refusée et les valeurs acceptées.

func TestInvalidStatusMessageNamesFieldValueAndAllowedSet(t *testing.T) {
	msg := invalidStatusMessage("banana")

	if !strings.Contains(msg, "status") {
		t.Errorf("le message ne nomme pas le champ : %q", msg)
	}
	if !strings.Contains(msg, "banana") {
		t.Errorf("le message ne rappelle pas la valeur refusée : %q", msg)
	}
	for _, allowed := range models.AllowedTaskStatuses {
		if !strings.Contains(msg, allowed) {
			t.Errorf("le message ne cite pas %q parmi les valeurs acceptées : %q", allowed, msg)
		}
	}
}

func TestInvalidPriorityMessageNamesFieldValueAndAllowedSet(t *testing.T) {
	msg := invalidPriorityMessage("banana")

	if !strings.Contains(msg, "priority") {
		t.Errorf("le message ne nomme pas le champ : %q", msg)
	}
	if !strings.Contains(msg, "banana") {
		t.Errorf("le message ne rappelle pas la valeur refusée : %q", msg)
	}
	for _, allowed := range models.AllowedTaskPriorities {
		if !strings.Contains(msg, allowed) {
			t.Errorf("le message ne cite pas %q parmi les valeurs acceptées : %q", allowed, msg)
		}
	}
}

// Les deux messages doivent rester distinguables : un client qui envoie les
// deux champs d'un coup doit savoir lequel a été refusé.
func TestInvalidMessagesDistinguishTheTwoFields(t *testing.T) {
	status := invalidStatusMessage("banana")
	priority := invalidPriorityMessage("banana")

	if status == priority {
		t.Fatalf("les deux messages sont identiques : %q", status)
	}
	if strings.Contains(status, "priority") {
		t.Errorf("le message de statut parle de priorité : %q", status)
	}
	if strings.Contains(priority, "status") {
		t.Errorf("le message de priorité parle de statut : %q", priority)
	}
}

// Le vocabulaire refusé côté serveur doit être exactement celui de models :
// ces cas figent la frontière acceptée/refusée telle que les handlers
// l'appliquent.
func TestValidationBoundaryMatchesTheVocabulary(t *testing.T) {
	statuses := map[string]bool{
		models.TaskStatusTodo:       true,
		models.TaskStatusInProgress: true,
		models.TaskStatusDone:       true,
		"DOING":                     false,
		"banana":                    false,
		"":                          false,
	}

	for status, want := range statuses {
		if got := models.IsValidTaskStatus(status); got != want {
			t.Errorf("IsValidTaskStatus(%q) = %v, attendu %v", status, got, want)
		}
	}

	priorities := map[string]bool{
		models.TaskPriorityLow:    true,
		models.TaskPriorityMedium: true,
		models.TaskPriorityHigh:   true,
		"URGENT":                  false,
		"":                        false,
	}

	for priority, want := range priorities {
		if got := models.IsValidTaskPriority(priority); got != want {
			t.Errorf("IsValidTaskPriority(%q) = %v, attendu %v", priority, got, want)
		}
	}
}
