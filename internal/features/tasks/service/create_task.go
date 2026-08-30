package tasks_service

import (
	"context"
	"fmt"

	"github.com/Muxammednone/golang-todoapp/internal/core/domain"
)

func (s *TasksService) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	//! провалидировать задачу task.Validate()
	//! newTask := repo.Save()
	//! return newTask

	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("validate task domain: %w", err)
	}

	task, err := s.tasksRepository.CreateTask(ctx, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}
