// Command seed remplit une base locale avec un jeu de donnees de test.
//
//	go run ./cmd/seed                          # sème, puis affiche les logins
//	go run ./cmd/seed -join moi@exemple.com    # rattache un compte existant
//	go run ./cmd/seed -force                   # passe outre le garde-fou d'hote
//
// Le semis est un binaire separe, et non un drapeau du serveur : il ne peut
// pas partir tout seul au demarrage d'un conteneur. Relancer remet le jeu de
// donnees a plat — les lignes semees sont effacees pour de bon, puis
// reconstruites. Ne sont concernees que les comptes en @example.test et les
// projets qu'ils possedent : le reste de la base n'est jamais touche.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"gorm.io/gorm"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
)

func main() {
	join := flag.String("join", "", "email d'un compte existant a ajouter comme MEMBER d'un projet seme")
	force := flag.Bool("force", false, "semer meme si DB_HOST n'est pas une base locale")
	flag.Parse()

	initializers.LoadEnvVariables()

	host := os.Getenv("DB_HOST")
	if !IsLocalDBHost(host) && !*force {
		log.Fatalf("refus de semer : DB_HOST=%q n'est pas une base locale.\n"+
			"Le semis efface des lignes pour de bon. Relancer avec -force pour passer outre.", host)
	}

	initializers.ConnectToDB()
	initializers.SyncDataBase()

	db := initializers.DB
	now := time.Now()

	removed, err := wipe(db)
	if err != nil {
		log.Fatalf("nettoyage impossible : %v", err)
	}
	if removed > 0 {
		fmt.Printf("\n%d compte(s) seme(s) au passage precedent ont ete effaces.\n", removed)
	}

	created, err := plant(db, now)
	if err != nil {
		log.Fatalf("semis impossible : %v", err)
	}

	if *join != "" {
		if err := joinSeededProject(db, *join); err != nil {
			log.Fatalf("rattachement de %s impossible : %v", *join, err)
		}
		fmt.Printf("\n%s ajoute comme MEMBER de %q.\n", *join, joinTargetName())
	}

	report(created)
}

// wipe efface les lignes semees au passage precedent. Unscoped, donc un vrai
// DELETE : l'index unique sur users.email ne fait pas de cas des suppressions
// douces de GORM, une ligne soft-deleted continuerait d'occuper l'adresse et le
// second passage echouerait sur un doublon.
func wipe(db *gorm.DB) (int64, error) {
	var userIDs []uint
	if err := db.Unscoped().Model(&models.User{}).
		Where("email LIKE ?", "%@"+SeedDomain).
		Pluck("id", &userIDs).Error; err != nil {
		return 0, err
	}
	if len(userIDs) == 0 {
		return 0, nil
	}

	var projectIDs []uint
	if err := db.Unscoped().Model(&models.Project{}).
		Where("owner_id IN ?", userIDs).
		Pluck("id", &projectIDs).Error; err != nil {
		return 0, err
	}

	var taskIDs []uint
	if len(projectIDs) > 0 {
		if err := db.Unscoped().Model(&models.Task{}).
			Where("project_id IN ?", projectIDs).
			Pluck("id", &taskIDs).Error; err != nil {
			return 0, err
		}
	}

	// Des feuilles vers la racine : les assignations, puis les taches, puis les
	// adhesions, puis les projets, et les comptes en dernier.
	if len(taskIDs) > 0 {
		if err := db.Unscoped().Where("task_id IN ?", taskIDs).Delete(&models.TaskAssignee{}).Error; err != nil {
			return 0, err
		}
		if err := db.Unscoped().Where("id IN ?", taskIDs).Delete(&models.Task{}).Error; err != nil {
			return 0, err
		}
	}
	// Une adhesion d'un compte seme peut viser un projet qui ne l'est pas (cas
	// -join a l'envers), et un projet seme peut heberger un compte qui ne l'est
	// pas : les deux cotes doivent partir.
	if err := db.Unscoped().Where("user_id IN ?", userIDs).Delete(&models.ProjectMember{}).Error; err != nil {
		return 0, err
	}
	if len(projectIDs) > 0 {
		if err := db.Unscoped().Where("project_id IN ?", projectIDs).Delete(&models.ProjectMember{}).Error; err != nil {
			return 0, err
		}
		if err := db.Unscoped().Where("id IN ?", projectIDs).Delete(&models.Project{}).Error; err != nil {
			return 0, err
		}
	}
	// Une tache creee par un compte seme dans un projet qui ne l'est pas
	// retiendrait la cle etrangere creator_id : on la laisse et on refuse
	// plutot que de toucher a un projet qui n'est pas a nous.
	if err := db.Unscoped().Where("id IN ?", userIDs).Delete(&models.User{}).Error; err != nil {
		return 0, err
	}

	return int64(len(userIDs)), nil
}

type counts struct {
	Users, Projects, Tasks, Assignments int
}

// plant insere le plan. Les cles du plan ("amina", "marc") sont resolues en
// identifiants au fur et a mesure : rien dans plan.go ne connait un id.
func plant(db *gorm.DB, now time.Time) (counts, error) {
	var made counts
	userIDs := map[string]uint{}

	for _, u := range SeedUsers {
		user := models.User{Name: u.Name, Email: u.Email}
		if err := user.SetPassword(SeedPassword); err != nil {
			return made, err
		}
		if err := db.Create(&user).Error; err != nil {
			return made, fmt.Errorf("compte %s : %w", u.Email, err)
		}
		userIDs[u.Key] = user.ID
		made.Users++
	}

	for _, p := range SeedProjects {
		ownerID := userIDs[p.OwnerKey]
		project := models.Project{Name: p.Name, Description: p.Description, OwnerID: &ownerID}
		if err := db.Create(&project).Error; err != nil {
			return made, fmt.Errorf("projet %s : %w", p.Name, err)
		}
		made.Projects++

		// Le proprietaire est aussi membre, avec le role OWNER : c'est ce que
		// fait CreateProject, et permission_helpers.go s'appuie sur les deux.
		members := []models.ProjectMember{{ProjectID: project.ID, UserID: ownerID, Role: models.RoleOwner}}
		for _, key := range p.MemberKeys {
			members = append(members, models.ProjectMember{
				ProjectID: project.ID, UserID: userIDs[key], Role: models.RoleMember,
			})
		}
		if err := db.Create(&members).Error; err != nil {
			return made, fmt.Errorf("membres de %s : %w", p.Name, err)
		}

		for _, t := range p.Tasks {
			task := models.Task{
				ProjectID:   project.ID,
				Title:       t.Title,
				Description: t.Description,
				Status:      t.Status,
				Priority:    t.Priority,
				DueDate:     DueAt(now, t.DueIn),
				CreatorID:   userIDs[t.CreatorKey],
			}
			if err := db.Create(&task).Error; err != nil {
				return made, fmt.Errorf("tache %s : %w", t.Title, err)
			}
			made.Tasks++

			for _, key := range t.AssigneeKeys {
				assignee := models.TaskAssignee{TaskID: task.ID, UserID: userIDs[key]}
				if err := db.Create(&assignee).Error; err != nil {
					return made, fmt.Errorf("assignation de %s : %w", t.Title, err)
				}
				made.Assignments++
			}
		}
	}

	return made, nil
}

// joinSeededProject ajoute un compte existant au projet marque JoinTarget. Il
// n'ecrit jamais dans un projet que le semis n'a pas cree.
func joinSeededProject(db *gorm.DB, email string) error {
	var user models.User
	if err := db.Where("email = ?", strings.ToLower(strings.TrimSpace(email))).First(&user).Error; err != nil {
		return fmt.Errorf("aucun compte avec cette adresse (%w)", err)
	}

	var project models.Project
	if err := db.Where("name = ?", joinTargetName()).First(&project).Error; err != nil {
		return err
	}

	return db.Create(&models.ProjectMember{
		ProjectID: project.ID, UserID: user.ID, Role: models.RoleMember,
	}).Error
}

func joinTargetName() string {
	for _, p := range SeedProjects {
		if p.JoinTarget {
			return p.Name
		}
	}
	return SeedProjects[0].Name
}

// report affiche la table des logins. Sans elle, le semis laisse cinq comptes
// dont personne ne connait l'adresse.
func report(made counts) {
	fmt.Printf("\nSeme : %d comptes, %d projets, %d taches, %d assignations.\n",
		made.Users, made.Projects, made.Tasks, made.Assignments)

	fmt.Printf("\nMot de passe pour tous les comptes ci-dessous : %s\n\n", SeedPassword)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "EMAIL\tNOM\tPROPRIETAIRE DE")
	for _, u := range SeedUsers {
		var owned []string
		for _, p := range SeedProjects {
			if p.OwnerKey == u.Key {
				owned = append(owned, p.Name)
			}
		}
		if owned == nil {
			owned = []string{"—"}
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", u.Email, u.Name, strings.Join(owned, ", "))
	}
	w.Flush()
}
