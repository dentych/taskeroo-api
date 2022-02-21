package app

import (
	"context"
	"github.com/dentych/taskeroo/internal/database"
	internalerrors "github.com/dentych/taskeroo/internal/errors"
	"github.com/google/uuid"
	"time"
)

type TaskLogic struct {
	taskRepo *database.TaskRepo
	userRepo *database.UserRepo
}

type Task struct {
	ID           string
	GroupID      string
	Title        string
	Description  string
	IntervalSize int
	IntervalUnit string
}

func NewTaskLogic(taskRepo *database.TaskRepo, userRepo *database.UserRepo) *TaskLogic {
	return &TaskLogic{taskRepo: taskRepo, userRepo: userRepo}
}

type NewTask struct {
	Title        string
	Description  string
	IntervalSize int
	IntervalUnit string
}

func (t *TaskLogic) Create(ctx context.Context, userID string, task NewTask) (Task, error) {
	user, err := t.userRepo.Get(ctx, userID)
	if err != nil {
		return Task{}, err
	}
	if user.GroupID == nil {
		return Task{}, internalerrors.ErrUserNotInGroup
	}

	taskID := uuid.NewString()
	err = t.taskRepo.Create(ctx, database.Task{
		ID:           taskID,
		Title:        task.Title,
		Description:  task.Description,
		GroupID:      *user.GroupID,
		IntervalSize: task.IntervalSize,
		IntervalUnit: task.IntervalUnit,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		return Task{}, err
	}

	return Task{
		ID:           taskID,
		GroupID:      *user.GroupID,
		Title:        task.Title,
		Description:  task.Description,
		IntervalSize: task.IntervalSize,
		IntervalUnit: task.IntervalUnit,
	}, nil
}

func (t *TaskLogic) GetForGroup(ctx context.Context, userID string, groupID string) ([]Task, error) {
	user, err := t.userRepo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.GroupID == nil || *user.GroupID != groupID {
		return nil, internalerrors.ErrUserNotMemberOfGroup
	}

	tasks, err := t.taskRepo.GetAllForGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	var mappedTasks []Task
	for _, task := range tasks {
		mappedTasks = append(mappedTasks, Task{
			ID:           task.ID,
			GroupID:      task.GroupID,
			Title:        task.Title,
			Description:  task.Description,
			IntervalSize: task.IntervalSize,
			IntervalUnit: task.IntervalUnit,
		})
	}

	return mappedTasks, nil
}

func (t *TaskLogic) Delete(ctx context.Context, userID string, taskID string) error {
	user, err := t.userRepo.Get(ctx, userID)
	if err != nil {
		return err
	}

	if user.GroupID == nil {
		return internalerrors.ErrUserNotInGroup
	}

	task, err := t.taskRepo.Get(ctx, taskID)
	if err != nil {
		return err
	}

	if task.GroupID != *user.GroupID {
		return internalerrors.ErrUserNotMemberOfGroup
	}

	err = t.taskRepo.Delete(ctx, taskID)
	if err != nil {
		return err
	}

	return nil
}
