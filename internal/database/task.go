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
	ID           string `gorm:"primaryKey;"`
	Title        string
	Description  string
	GroupID      string `gorm:"index"`
	IntervalSize int
	IntervalUnit string
	CreatedAt    time.Time
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
