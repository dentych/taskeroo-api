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

func NewTeamRepo(db *gorm.DB) *GroupRepo {
	return &GroupRepo{db: db}
}

func (r *GroupRepo) Create(ctx context.Context, group Group) error {
	return r.db.WithContext(ctx).Create(&group).Error
}
