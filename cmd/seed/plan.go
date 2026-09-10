package main

import (
	"strings"
	"time"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
)

// Le jeu de donnees est decrit ici, sans base : ce fichier ne fait que produire
// un plan. C'est ce qui le rend testable, et c'est aussi ce qui le rend
// lisible — la couverture du vocabulaire de CONTEXT.md se verifie a l'oeil.

// SeedDomain : les comptes semes vivent tous sur ce domaine. `.test` est un TLD
// reserve (RFC 2606) : aucune de ces adresses ne peut etre une vraie boite, et
// le domaine sert d'etiquette pour savoir quoi effacer au passage suivant.
const SeedDomain = "example.test"

// SeedPassword : le meme mot de passe pour tous les comptes semes. Ces comptes
// n'existent que sur un MySQL local jetable ; leur devinabilite est le but.
const SeedPassword = "secret123"

type userPlan struct {
	Key   string // reference interne au plan, jamais stockee
	Name  string
	Email string
}

type taskPlan struct {
	Title       string
	Description string
	Status      string
	Priority    string
	// DueIn est l'ecart a l'instant du semis. nil = pas d'echeance, donc aucun
	// Reminder possible. Un ecart plutot qu'une date : le jeu de donnees garde
	// le meme sens dans un mois.
	DueIn        *time.Duration
	CreatorKey   string
	AssigneeKeys []string
}

type projectPlan struct {
	Name        string
	Description string
	OwnerKey    string
	MemberKeys  []string
	Tasks       []taskPlan
	// JoinTarget marque le projet auquel -join rattache le compte donne.
	JoinTarget bool
}

func in(d time.Duration) *time.Duration { return &d }

const (
	hour = time.Hour
	day  = 24 * time.Hour
)

// SeedUsers : le casting. Cinq comptes suffisent a couvrir les deux roles et
// l'assignation multiple sans rendre les listes illisibles.
var SeedUsers = []userPlan{
	{Key: "amina", Name: "Amina Chakroun", Email: "amina@" + SeedDomain},
	{Key: "marc", Name: "Marc Lefevre", Email: "marc@" + SeedDomain},
	{Key: "sofia", Name: "Sofia Rossi", Email: "sofia@" + SeedDomain},
	{Key: "tomas", Name: "Tomas Novak", Email: "tomas@" + SeedDomain},
	{Key: "lena", Name: "Lena Haddad", Email: "lena@" + SeedDomain},
}

// SeedProjects : chaque terme de CONTEXT.md a ici au moins une instance
// visible. Les cas de contraste — une tache DONE dont la date est passee, une
// tache sans assigne et sans echeance — sont la expres : ce sont eux qui
// distinguent Overdue de « en retard » et Orphan de « non assignee ».
var SeedProjects = []projectPlan{
	{
		Name:        "Claims Triage Agent",
		Description: "Route incoming claims to the right adjuster, and flag the ones a human must read first.",
		OwnerKey:    "amina",
		MemberKeys:  []string{"marc", "sofia"},
		JoinTarget:  true,
		Tasks: []taskPlan{
			{
				Title:        "Draft the claim intake schema",
				Description:  "The fields every claim carries, whatever channel it arrives on.",
				Status:       models.TaskStatusDone,
				Priority:     models.TaskPriorityMedium,
				DueIn:        in(-10 * day), // DONE et depassee : ne doit PAS ressortir en Overdue
				CreatorKey:   "amina",
				AssigneeKeys: []string{"marc"},
			},
			{
				Title:        "Wire OCR for scanned claim forms",
				Description:  "Half the intake still arrives as a photograph of a paper form.",
				Status:       models.TaskStatusInProgress,
				Priority:     models.TaskPriorityHigh,
				DueIn:        in(-3 * day), // Overdue
				CreatorKey:   "amina",
				AssigneeKeys: []string{"marc", "sofia"}, // deux assignes
			},
			{
				Title:        "Escalation rules for high-value claims",
				Description:  "Anything above the threshold goes to a human, no exceptions.",
				Status:       models.TaskStatusTodo,
				Priority:     models.TaskPriorityHigh,
				DueIn:        in(6 * hour), // Due soon
				CreatorKey:   "amina",
				AssigneeKeys: []string{"amina"},
			},
			{
				Title:       "Redact PII before the model call",
				Description: "Names, policy numbers and addresses must never leave the perimeter.",
				Status:      models.TaskStatusTodo,
				Priority:    models.TaskPriorityHigh,
				DueIn:       in(-1 * day), // Orphan : porte un Reminder, personne ne la doit
				CreatorKey:  "amina",
			},
			{
				Title:        "Benchmark latency against the manual baseline",
				Description:  "The agent has to beat an adjuster reading the file, not merely be correct.",
				Status:       models.TaskStatusTodo,
				Priority:     models.TaskPriorityLow,
				DueIn:        in(21 * day), // ni Overdue ni Due soon
				CreatorKey:   "sofia",
				AssigneeKeys: []string{"sofia"},
			},
			{
				Title:        "Write the evaluation rubric",
				Description:  "What counts as a good triage decision, written down before we measure any.",
				Status:       models.TaskStatusInProgress,
				Priority:     models.TaskPriorityMedium,
				CreatorKey:   "amina", // sans echeance : aucun Reminder possible
				AssigneeKeys: []string{"amina"},
			},
			{
				Title:        "Retire the legacy rules engine",
				Description:  "Kept running in parallel until the agent has a month of clean decisions.",
				Status:       models.TaskStatusDone,
				Priority:     models.TaskPriorityLow,
				DueIn:        in(14 * day),
				CreatorKey:   "marc",
				AssigneeKeys: []string{"marc"},
			},
		},
	},
	{
		Name:        "Policy Renewal Assistant",
		Description: "Draft renewal offers and chase the policies about to lapse.",
		OwnerKey:    "marc",
		MemberKeys:  []string{"amina", "tomas"},
		Tasks: []taskPlan{
			{
				Title:        "Import the renewal calendar",
				Description:  "Twelve months of expiry dates, one CSV per branch.",
				Status:       models.TaskStatusDone,
				Priority:     models.TaskPriorityHigh,
				DueIn:        in(-20 * day),
				CreatorKey:   "marc",
				AssigneeKeys: []string{"marc"},
			},
			{
				Title:        "Draft the renewal reminder copy",
				Description:  "One version for motor, one for home. Legal reviews both.",
				Status:       models.TaskStatusTodo,
				Priority:     models.TaskPriorityMedium,
				DueIn:        in(2 * day),
				CreatorKey:   "marc",
				AssigneeKeys: []string{"tomas"},
			},
			{
				Title:        "Handle mid-term policy changes",
				Description:  "A policy amended in month seven must not be renewed on its old terms.",
				Status:       models.TaskStatusInProgress,
				Priority:     models.TaskPriorityHigh,
				DueIn:        in(-2 * hour), // Overdue de peu
				CreatorKey:   "marc",
				AssigneeKeys: []string{"amina"},
			},
			{
				Title:       "Pricing edge cases nobody owns yet",
				Description: "Multi-vehicle, mid-term cancellation, and the two legacy tariffs.",
				Status:      models.TaskStatusTodo,
				Priority:    models.TaskPriorityLow,
				DueIn:       in(18 * hour), // Orphan, cote Due soon cette fois
				CreatorKey:  "marc",
			},
			{
				Title:       "Localise the renewal emails to French",
				Description: "Same content, reviewed by the Tunis office.",
				Status:      models.TaskStatusTodo,
				Priority:    models.TaskPriorityLow,
				CreatorKey:  "marc", // ni echeance ni assigne : non assignee, mais pas Orphan
			},
		},
	},
	{
		Name:        "Fraud Signals Backlog",
		Description: "Signals worth investigating, once somebody has the time to write them up.",
		OwnerKey:    "sofia",
		MemberKeys:  []string{"lena"},
		// Aucune tache : c'est l'etat vide de la page projet.
	},
}

// IsLocalDBHost dit si l'hote donne est une base locale jetable. Le semis
// efface pour de bon : il ne doit pas etre facile de le pointer ailleurs que
// sur un poste de developpement. "mysql" est le nom du service dans
// docker-compose.yml.
func IsLocalDBHost(host string) bool {
	switch strings.ToLower(strings.TrimSpace(host)) {
	case "", "localhost", "127.0.0.1", "::1", "[::1]", "mysql":
		return true
	}
	return false
}

// DueAt convertit l'ecart d'une tache en date absolue, relativement a
// l'instant du semis.
func DueAt(now time.Time, dueIn *time.Duration) *time.Time {
	if dueIn == nil {
		return nil
	}
	at := now.Add(*dueIn)
	return &at
}
