package database

import (
	"context"
	"gorm.io/gorm"
)

type GroupRepo struct {
	db *gorm.DB
}

type Group struct {
	ID          string `gorm:"primaryKey;"`
	Name        string
	OwnerUserID string
}

func NewGroupRepo(db *gorm.DB) *GroupRepo {
	return &GroupRepo{db: db}
}

func (r *GroupRepo) Create(ctx context.Context, group Group) error {
	return r.db.WithContext(ctx).Create(&group).Error
}

func (r *GroupRepo) Get(ctx context.Context, groupID string) (*Group, error) {
	var group Group
	err := r.db.WithContext(ctx).First(&group, "id = ?", groupID).Error
	if err != nil {
		return nil, err
	}

	return &group, nil
}
