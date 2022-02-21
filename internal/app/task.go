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
	ID             string
	GroupID        string
	Title          string
	Description    string
	IntervalSize   int
	IntervalUnit   string
	DaysLeft       int
	PercentageLeft float64
	DueDate        string
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

func (t *TaskLogic) Create(ctx context.Context, userID string, newTask NewTask) (Task, error) {
	user, err := t.userRepo.Get(ctx, userID)
	if err != nil {
		return Task{}, err
	}
	if user.GroupID == nil {
		return Task{}, internalerrors.ErrUserNotInGroup
	}

	taskID := uuid.NewString()
	task := database.Task{
		ID:           taskID,
		Title:        newTask.Title,
		Description:  newTask.Description,
		GroupID:      *user.GroupID,
		IntervalSize: newTask.IntervalSize,
		IntervalUnit: newTask.IntervalUnit,
		NextDueDate:  calculateNextDueDate(newTask.IntervalUnit, newTask.IntervalSize),
		CreatedAt:    time.Now(),
	}
	err = t.taskRepo.Create(ctx, task)
	if err != nil {
		return Task{}, err
	}

	return Task{
		ID:             taskID,
		GroupID:        *user.GroupID,
		Title:          newTask.Title,
		Description:    newTask.Description,
		IntervalSize:   newTask.IntervalSize,
		IntervalUnit:   newTask.IntervalUnit,
		DaysLeft:       calculateDaysLeft(task.NextDueDate),
		PercentageLeft: calculatePercentageLeft(task.IntervalUnit, task.IntervalSize, task.NextDueDate),
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
			ID:             task.ID,
			GroupID:        task.GroupID,
			Title:          task.Title,
			Description:    task.Description,
			IntervalSize:   task.IntervalSize,
			IntervalUnit:   task.IntervalUnit,
			DaysLeft:       calculateDaysLeft(task.NextDueDate),
			PercentageLeft: calculatePercentageLeft(task.IntervalUnit, task.IntervalSize, task.NextDueDate),
			DueDate:        dateFormat(task.NextDueDate),
		})
	}

	return mappedTasks, nil
}

func dateFormat(date time.Time) string {
	return date.Format("02/01/2006")
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

func calculateDaysLeft(date time.Time) int {
	until := time.Until(date)
	if until.Hours() < 24 {
		if until.Hours() < 12 {
			return 0
		}
		return 1
	}
	return int(until.Truncate(24*time.Hour).Hours()) / 24
}

func calculatePercentageLeft(unit string, size int, nextDueDate time.Time) float64 {
	totalHours := calculateTotalHours(unit, size)

	hoursUntilDue := time.Until(nextDueDate).Hours()

	result := hoursUntilDue / float64(totalHours)
	if result < 0 {
		return 0
	}

	return result
}

func calculateTotalHours(unit string, size int) int {
	switch unit {
	case "day":
		return 24 * size
	case "week":
		return 24 * 7 * size
	case "month":
		return 24 * 30 * size
	default:
		return 24
	}
}

func calculateNextDueDate(unit string, size int) time.Time {
	switch unit {
	case "day":
		return time.Now().AddDate(0, 0, size)
	case "week":
		return time.Now().AddDate(0, 0, size*7)
	case "month":
		return time.Now().AddDate(0, size, 0)
	default:
		return time.Now()
	}
}

var day = map[time.Weekday]string{
	time.Monday:    "Mandag",
	time.Tuesday:   "Tirsdag",
	time.Wednesday: "Onsdag",
	time.Thursday:  "Torsdag",
	time.Friday:    "Fredag",
	time.Saturday:  "Lørdag",
	time.Sunday:    "Søndag",
}

var month = map[time.Month]string{
	time.January:   "Januar",
	time.February:  "Februar",
	time.March:     "Marts",
	time.April:     "April",
	time.May:       "Maj",
	time.June:      "Juni",
	time.July:      "Juli",
	time.August:    "August",
	time.September: "September",
	time.October:   "Oktober",
	time.November:  "November",
	time.December:  "December",
}
