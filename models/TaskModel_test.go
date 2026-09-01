package models

import "testing"

// Le vocabulaire des tâches est partagé avec le frontend et avec les lignes
// déjà stockées en base : ces tests le figent des deux côtés à la fois, en
// vérifiant les valeurs elles-mêmes et pas seulement les constantes.

func TestTaskStatusConstantsMatchStoredValues(t *testing.T) {
	cases := map[string]string{
		TaskStatusTodo:       "TODO",
		TaskStatusInProgress: "IN_PROGRESS",
		TaskStatusDone:       "DONE",
	}

	for got, want := range cases {
		if got != want {
			t.Errorf("statut = %q, attendu %q", got, want)
		}
	}
}

func TestTaskPriorityConstantsMatchStoredValues(t *testing.T) {
	cases := map[string]string{
		TaskPriorityLow:    "LOW",
		TaskPriorityMedium: "MEDIUM",
		TaskPriorityHigh:   "HIGH",
	}

	for got, want := range cases {
		if got != want {
			t.Errorf("priorité = %q, attendu %q", got, want)
		}
	}
}

func TestIsValidTaskStatus(t *testing.T) {
	cases := []struct {
		status string
		want   bool
	}{
		{TaskStatusTodo, true},
		{TaskStatusInProgress, true},
		{TaskStatusDone, true},
		{"banana", false},
		{"", false},
		{"DOING", false},         // l'ancienne orthographe n'est plus acceptée
		{"in_progress", false},   // la casse compte
		{" IN_PROGRESS ", false}, // rien n'est rogné avant comparaison
	}

	for _, tc := range cases {
		if got := IsValidTaskStatus(tc.status); got != tc.want {
			t.Errorf("IsValidTaskStatus(%q) = %v, attendu %v", tc.status, got, tc.want)
		}
	}
}

func TestIsValidTaskPriority(t *testing.T) {
	cases := []struct {
		priority string
		want     bool
	}{
		{TaskPriorityLow, true},
		{TaskPriorityMedium, true},
		{TaskPriorityHigh, true},
		{"URGENT", false},
		{"", false},
		{"low", false},
	}

	for _, tc := range cases {
		if got := IsValidTaskPriority(tc.priority); got != tc.want {
			t.Errorf("IsValidTaskPriority(%q) = %v, attendu %v", tc.priority, got, tc.want)
		}
	}
}

func TestAllowedSetsCoverEveryConstant(t *testing.T) {
	if len(AllowedTaskStatuses) != 3 {
		t.Errorf("AllowedTaskStatuses contient %d valeurs, attendu 3", len(AllowedTaskStatuses))
	}
	if len(AllowedTaskPriorities) != 3 {
		t.Errorf("AllowedTaskPriorities contient %d valeurs, attendu 3", len(AllowedTaskPriorities))
	}
}
