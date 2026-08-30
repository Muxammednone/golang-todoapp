package tasks_service

import (
	"context"
	"fmt"

	"github.com/Muxammednone/golang-todoapp/internal/core/domain"
)

func (s *TasksService) PatchTask(
	ctx context.Context,
	taskID int,
	taskPatch domain.TaskPatch,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}

	if err := task.ApplyPatch(taskPatch); err != nil {
		return domain.Task{}, fmt.Errorf("apply task patch: %w", err)
	}

	patchedTask, err := s.tasksRepository.PatchTask(ctx, taskID, task)

	if err != nil {
		return domain.Task{}, fmt.Errorf("patch task: %w", err)
	}

	return patchedTask, nil
}
