package tasks_postgres_repository

import (
	"time"

	"github.com/Muxammednone/golang-todoapp/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func taskDomainsFromModels(taskModels []TaskModel) []domain.Task {
	taskDomains := make([]domain.Task, len(taskModels))

	for i, taskModel := range taskModels {
		taskDomains[i] = taskDomainFromModel(taskModel)
	}

	return taskDomains
}

func taskDomainFromModel(taskModel TaskModel) domain.Task {
	return domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreatedAt,
		taskModel.CompletedAt,
		taskModel.AuthorUserID,
	)
}
