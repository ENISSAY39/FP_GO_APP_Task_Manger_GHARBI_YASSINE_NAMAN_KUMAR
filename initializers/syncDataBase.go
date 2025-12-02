package initializers

import (
	"fmt"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
)

func SyncDataBase() {
	if DB == nil {
		fmt.Println("SyncDataBase: DB is nil")
		return
	}
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.ProjectMember{},
		&models.Task{},
		&models.TaskAssignee{},
	); err != nil {
		fmt.Println("AutoMigrate error:", err)
	} else {
		fmt.Println("AutoMigrate completed")
	}
}
