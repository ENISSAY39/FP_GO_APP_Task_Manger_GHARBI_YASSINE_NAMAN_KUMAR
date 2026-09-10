package controllers

import (
	"net/http"
	"strings"
	"testing"
)

// resolve() ne touche la base que lorsqu'un email est fourni. Les branches
// sans base — id seul, et payload vide — sont donc testables telles quelles,
// et ce sont elles qui portent le contrat vu du client.

func TestResolveAcceptsAUserIDWithoutTouchingTheDatabase(t *testing.T) {
	id, msg, status := addMemberPayload{UserID: 7}.resolve()

	if msg != "" {
		t.Fatalf("un id seul devrait suffire, refus : %q (status %d)", msg, status)
	}
	if id != 7 {
		t.Errorf("id résolu %d, attendu 7", id)
	}
}

func TestResolveRefusesAnEmptyPayloadByNamingTheEmail(t *testing.T) {
	_, msg, status := addMemberPayload{}.resolve()

	if msg == "" {
		t.Fatal("un payload vide devrait être refusé")
	}
	if !strings.Contains(msg, "email") {
		t.Errorf("le message n'oriente pas vers l'email : %q", msg)
	}
	if status != http.StatusBadRequest {
		t.Errorf("status %d, attendu %d", status, http.StatusBadRequest)
	}
}

// Un email réduit à des espaces n'identifie personne : il ne doit pas être
// envoyé à la base comme s'il était une adresse.
func TestResolveTreatsABlankEmailAsAbsent(t *testing.T) {
	_, msg, status := addMemberPayload{Email: "   "}.resolve()

	if status != http.StatusBadRequest {
		t.Errorf("status %d, attendu %d (message : %q)", status, http.StatusBadRequest, msg)
	}
}
