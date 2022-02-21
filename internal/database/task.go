package database

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type TaskRepo struct {
	db *gorm.DB
}

type Task struct {
	ID               string `gorm:"primaryKey;"`
	Title            string `gorm:"not null;"`
	Description      string
	GroupID          string  `gorm:"index"`
	Assignee         *string `gorm:"index"`
	RotatingAssignee bool    `gorm:"not null;default: false;"`
	IntervalSize     int     `gorm:"not null;"`
	IntervalUnit     string  `gorm:"not null;"`
	NextDueDate      time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

func NewTaskRepo(db *gorm.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(ctx context.Context, task Task) error {
	return r.db.WithContext(ctx).Create(&task).Error
}

func (r *TaskRepo) GetAllForGroup(ctx context.Context, groupID string) ([]Task, error) {
	var tasks []Task
	err := r.db.WithContext(ctx).Find(&tasks, "group_id = ?", groupID).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepo) Delete(ctx context.Context, taskID string) error {
	return r.db.WithContext(ctx).Delete(&Task{ID: taskID}).Error
}

func (r *TaskRepo) Get(ctx context.Context, taskID string) (*Task, error) {
	var task Task
	err := r.db.WithContext(ctx).First(&task, "id = ?", taskID).Error
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepo) Update(ctx context.Context, task Task) error {
	return r.db.WithContext(ctx).Model(&task).Updates(map[string]interface{}{
		"title":             task.Title,
		"description":       task.Description,
		"group_id":          task.GroupID,
		"assignee":          task.Assignee,
		"rotating_assignee": task.RotatingAssignee,
		"interval_size":     task.IntervalSize,
		"interval_unit":     task.IntervalUnit,
		"next_due_date":     task.NextDueDate,
		"updated_at":        time.Now(),
	}).Error
}
